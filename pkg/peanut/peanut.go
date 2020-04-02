package peanut

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type FileMapping struct {
	selector FileSelector
	destPath string
}

type RunConfiguration struct {
	Sources []Source
}

type Source struct {
	fs       SourceFilesystem
	mappings []FileMapping
}

func NewFileMapping(selector FileSelector, destPath string) FileMapping {
	return FileMapping{selector, destPath}
}

func NewSource(fs SourceFilesystem, mappings []FileMapping) *Source {
	return &Source{fs, mappings}
}

// Copies the files from the source to the destination directories
func (source *Source) Pull(destDir string) error {
	dir, err := ioutil.TempDir("", "peanut")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	err = source.fs.FetchAll(dir)
	if err != nil {
		return err
	}

	for _, mapping := range source.mappings {
		files, err := mapping.selector.GetFiles(dir)
		if err != nil {
			return err
		}

		for _, file := range files {
			fullDestPath := path.Join(destDir, mapping.destPath)

				err = file.CopyTo(filepath.FromSlash(fullDestPath))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
