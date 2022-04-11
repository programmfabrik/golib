package golib

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// UnpackZipFile reads zipData and unpacks the contents of the zip
// into targetDir. targetDir must exist.
func UnpackZipFile(targetDir string, zipData io.Reader) (err error) {
	zipFile, err := ioutil.TempFile(targetDir, "")
	if err != nil {
		return errors.Wrap(err, "Error installing plugin. Writing ZIP failed")
	}

	defer func() {
		zipFile.Close()
		_ = os.Remove(zipFile.Name())
	}()

	_, err = io.Copy(zipFile, zipData)
	if err != nil {
		return errors.Wrap(err, "Error installing plugin. Copying ZIP failed")
	}

	r, err := zip.OpenReader(zipFile.Name())
	if err != nil {
		return errors.Wrap(err, "Opening ZIP failed")
	}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {

		fn := filepath.Join(targetDir, f.Name)

		if f.FileInfo().IsDir() {
			// Make Folder
			err = os.MkdirAll(fn, 0755)
			if err != nil {
				return errors.Wrapf(err, "Unpacking ZIP failed: %q", f.Name)
			}
			continue
		}

		// Make directory for file
		err = os.MkdirAll(filepath.Dir(fn), 0755)
		if err != nil {
			return errors.Wrapf(err, "Unpacking ZIP failed: %q", f.Name)
		}

		// Open file in ZIP
		rc, err := f.Open()
		if err != nil {
			return errors.Wrapf(err, "Unpacking ZIP failed: %q", f.Name)
		}

		// Open file on disk
		of, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return errors.Wrapf(err, "Unpacking ZIP failed: %q", f.Name)
		}

		// Copy data from ZIP to disk
		_, err = io.Copy(of, rc)
		if err != nil {
			of.Close()
			return errors.Wrapf(err, "Copy to file failed: %q", fn)
		}
		of.Close()
		rc.Close()
	}
	r.Close()
	return nil
}

// PackZipFile packs all files from sourceDir as ZIP to writeTo.
func PackZipFile(sourceDir, topLevelDir string, writeTo io.Writer) (err error) {
	archive := zip.NewWriter(writeTo)
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fn := strings.TrimPrefix(path, sourceDir+"/")
		if topLevelDir != "" {
			fn = topLevelDir + "/" + fn
		}
		fw, err := archive.Create(fn)
		if err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, file)
		file.Close()
		return err
	})
	archive.Close()
	if err != nil {
		return errors.Wrap(err, "Error packing ZIP")
	}
	return nil
}
