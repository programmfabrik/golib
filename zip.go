package golib

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// UnpackZipFile reads zipData and unpacks the contents of the zip
// into targetDir. targetDir must exist.
func UnpackZipFile(targetDir string, zipData io.Reader) (err error) {
	zipFile, err := os.CreateTemp(targetDir, "")
	if err != nil {
		return fmt.Errorf("Unpack ZIP: %w", err)
	}

	defer func() {
		zipFile.Close()
		_ = os.Remove(zipFile.Name())
	}()

	_, err = io.Copy(zipFile, zipData)
	if err != nil {
		return fmt.Errorf("Copy ZIP: %w", err)
	}

	r, err := zip.OpenReader(zipFile.Name())
	if err != nil {
		return fmt.Errorf("Read ZIP: %w", err)
	}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {

		// Join turns "/" into "\" on Windows
		fn := filepath.Join(targetDir, f.Name)

		if f.FileInfo().IsDir() {
			// Make Folder
			err = os.MkdirAll(fn, f.Mode())
			if err != nil {
				return fmt.Errorf("Mkdir for ZIP %q: %w", f.Name, err)
			}
			continue
		}

		// Make directory for file
		err = os.MkdirAll(filepath.Dir(fn), 0755)
		if err != nil {
			return fmt.Errorf("Mkdir for ZIP %q: %w", f.Name, err)
		}

		// Open file in ZIP
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("Open for ZIP %q: %w", f.Name, err)
		}

		// Open file on disk
		of, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("Open on disk for ZIP %q: %w", f.Name, err)
		}

		// Copy data from ZIP to disk
		_, err = io.Copy(of, rc)
		if err != nil {
			of.Close()
			return fmt.Errorf("Copy for ZIP %q: %w", fn, err)
		}
		of.Close()
		rc.Close()
	}
	r.Close()
	return nil
}

// PackZipFile packs all files from sourceDir as ZIP to writeTo. We pack
// all filenames using /. On Windows, the \ in paths is replaced by /
func PackZipFile(sourceDir, topLevelDir string, writeTo io.Writer) (err error) {
	sep := string(os.PathSeparator)
	zipW := zip.NewWriter(writeTo)
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
		fh, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		fh.Name = fn
		fw, err := zipW.CreateHeader(fh)
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
	zipW.Close()
	if err != nil {
		return fmt.Errorf("Error packing ZIP: %w", err)
	}
	return nil
}
