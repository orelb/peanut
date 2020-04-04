package peanut

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// Source describes how to fetch remote files and order them locally.
type Source struct {
	fs       SourceFilesystem
	mappings []FileMapping
}

// NewSource creates a new source
func NewSource(fs SourceFilesystem, mappings []FileMapping) *Source {
	return &Source{fs, mappings}
}

// Pull copies the files from the source to the destination directories
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
		files, err := matchFiles(dir, mapping.MatchPattern)
		if err != nil {
			return err
		}

		for _, file := range files {
			fullDestPath := path.Join(destDir, mapping.DestPath)
			fsPath := filepath.FromSlash(fullDestPath)

			err = file.CopyTo(fsPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
