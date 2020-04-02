package peanut

import (
	"io/ioutil"
	"path"
	"path/filepath"
)

type FileSelector interface {
	// Returns filenames selected by this selector
	GetFiles(sourceDirectory string) ([]LocalAwareFile, error)
}

type standardFileSelector struct {
	path string
}

func NewStandardFileSelector(path string) FileSelector {
	return standardFileSelector{path}
}

func (selector standardFileSelector) GetFiles(sourceDirectory string) ([]LocalAwareFile, error) {
	fullPath := path.Join(sourceDirectory, selector.path)
	fsPath := filepath.FromSlash(fullPath)

	fileInfo, err := AppFs.Stat(fsPath)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(fullPath)

		if err != nil {
			return nil, err
		}

		// Return only filenames
		laFiles := make([]LocalAwareFile, len(files))
		for i, file := range files {
			laFiles[i] = newLocalAwareFile(fullPath, file.Name())
		}

		return laFiles, nil
	} else {
		laFile := newLocalAwareFile(sourceDirectory, fileInfo.Name())
		return []LocalAwareFile{laFile}, nil
	}
}

type globFileSelector struct {
	pattern string
}

func NewGlobFileSelector(pattern string) FileSelector {
	return globFileSelector{pattern}
}

func (selector globFileSelector) GetFiles(sourceDirectory string) ([]LocalAwareFile, error) {
	fullPath := path.Join(sourceDirectory, selector.pattern)
	fsPath := filepath.FromSlash(fullPath)

	matches, err := filepath.Glob(fsPath)
	if err != nil {
		return nil, err
	}

	laFiles := make([]LocalAwareFile, len(matches))

	for i, matchPath := range matches {
		relativeFsPath, err := filepath.Rel(sourceDirectory, matchPath)
		if err != nil {
			return nil, err
		}

		relativePath := filepath.ToSlash(relativeFsPath)
		laFiles[i] = newLocalAwareFile(sourceDirectory, relativePath)
	}
	return laFiles, nil
}
