package golib

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// UnpackZipFile reads zipData and unpacks the contents of the zip
// into targetDir. targetDir must exist.
func UnpackZipFile(targetDir string, zipData io.Reader) (err error) {
	zipFile, err := ioutil.TempFile(targetDir, "")
	if err != nil {
		return errors.Wrap(err, "Unpack ZIP")
	}

	defer func() {
		zipFile.Close()
		_ = os.Remove(zipFile.Name())
	}()

	_, err = io.Copy(zipFile, zipData)
	if err != nil {
		return errors.Wrap(err, "Copy ZIP")
	}

	r, err := zip.OpenReader(zipFile.Name())
	if err != nil {
		return errors.Wrap(err, "Read ZIP")
	}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {

		// Join turns "/" into "\" on Windows
		fn := filepath.Join(targetDir, f.Name)

		if f.FileInfo().IsDir() {
			// Make Folder
			err = os.MkdirAll(fn, 0755)
			if err != nil {
				return errors.Wrapf(err, "Mkdir for ZIP %q", f.Name)
			}
			continue
		}

		// Make directory for file
		err = os.MkdirAll(filepath.Dir(fn), 0755)
		if err != nil {
			return errors.Wrapf(err, "Mkdir for ZIP %q", f.Name)
		}

		// Open file in ZIP
		rc, err := f.Open()
		if err != nil {
			return errors.Wrapf(err, "Open for ZIP %q", f.Name)
		}

		// Open file on disk
		of, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return errors.Wrapf(err, "Open on disk for ZIP %q", f.Name)
		}

		// Copy data from ZIP to disk
		_, err = io.Copy(of, rc)
		if err != nil {
			of.Close()
			return errors.Wrapf(err, "Copy for ZIP %q", fn)
		}
		of.Close()
		rc.Close()
	}
	r.Close()
	return nil
}

// PackZipFile packs all files from sourceDir as ZIP to writeTo. We pack
// all filenames using /. On Windows, the \ in paths is replaced by \
func PackZipFile(sourceDir, topLevelDir string, writeTo io.Writer) (err error) {
	sep := string(os.PathSeparator)
	archive := zip.NewWriter(writeTo)
	err = filepath.Walk(sourceDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fn := strings.ReplaceAll(strings.TrimPrefix(p, sourceDir+sep), sep, "/")
		if topLevelDir != "" {
			fn = path.Join(topLevelDir, fn) // Use "/" join here, independent from the OS as it is standard in ZIP
		}
		fw, err := archive.Create(fn)
		if err != nil {
			return err
		}
		file, err := os.Open(p)
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
