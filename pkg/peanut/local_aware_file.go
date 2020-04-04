package peanut

import (
	"os"
	"path"
	"path/filepath"
)

// localAwareFile remembers the baseDir it was matched on by matchFiles()
// It stores the paths in forward-slashes regardless to OS.
type localAwareFile struct {
	basePath     string
	relativePath string
}

func newLocalAwareFile(basePath, relativePath string) localAwareFile {
	return localAwareFile{filepath.ToSlash(basePath), filepath.ToSlash(relativePath)}
}

func (laFile *localAwareFile) Path() string {
	return path.Join(laFile.basePath, laFile.relativePath)
}

func (laFile *localAwareFile) FsPath() string {
	return filepath.FromSlash(laFile.Path())
}

func (laFile *localAwareFile) BasePath() string {
	return laFile.basePath
}

func (laFile *localAwareFile) RelativePath() string {
	return laFile.relativePath
}

func (laFile *localAwareFile) CopyTo(destDir string) error {
	// TODO: add support for copying directories
	fileInfo, err := fs.Stat(laFile.FsPath())
	if err != nil {
		return err
	}

	fullDestPath := filepath.Join(destDir, laFile.relativePath)
	fullDestDir := filepath.Dir(fullDestPath)

	err = fs.MkdirAll(fullDestDir, os.ModePerm)
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
