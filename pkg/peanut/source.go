package peanut

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type Source struct {
	fs       SourceFilesystem
	mappings []FileMapping
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
		files, err := MatchFiles(dir, mapping.MatchPattern)
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
