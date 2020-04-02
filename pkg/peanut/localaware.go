package peanut

import (
	"os"
	"path"
	"path/filepath"
)

type LocalAwareFile struct {
	basePath     string
	relativePath string
}

func newLocalAwareFile(basePath, relativePath string) LocalAwareFile {
	return LocalAwareFile{basePath, relativePath}
}

func (laFile *LocalAwareFile) Path() string {
	return path.Join(laFile.basePath, laFile.relativePath)
}

func (laFile *LocalAwareFile) FsPath() string {
	return filepath.FromSlash(laFile.Path())
}

func (laFile *LocalAwareFile) CopyTo(destDir string) error {
	fileInfo, err := AppFs.Stat(laFile.FsPath())
	if err != nil {
		return err
	}

	fullDestPath := filepath.Join(destDir, laFile.relativePath)
	fullDestDir := filepath.Dir(fullDestPath)

	err = AppFs.MkdirAll(fullDestDir, os.ModePerm)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		err = copyFile(laFile.FsPath(), fullDestPath)
		if err != nil {
			return err
		}
	}

	return nil
}
