package peanut

import (
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
)

func copyDirectory(sourceDir string, destinationDir string) error {
	err := afero.Walk(fs, sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(destinationDir, relativePath)

		if info.IsDir() {
			err = fs.MkdirAll(targetPath, os.ModePerm)

			if err != nil {
				return err
			}
		} else {
			err = copyFile(path, targetPath)

			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func copyFile(sourcePath string, destinationPath string) error {
	destination, err := fs.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	source, err := fs.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}
