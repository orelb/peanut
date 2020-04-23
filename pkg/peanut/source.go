package peanut

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// Source describes how to fetch remote files and order them locally.
type Source struct {
	name     string
	fs       SourceFilesystem
	mappings []FileMapping
}

// NewSource creates a new source
func NewSource(name string, fs SourceFilesystem, mappings []FileMapping) *Source {
	return &Source{name, fs, mappings}
}

// Name returns the source name.
func (source *Source) Name() string {
	return source.name
}

// Fetch copies the files from the source to the destination directories
func (source *Source) Fetch(destDir string) error {
	// Create a temporary directory, we first fetch the source's files to there.
	dir, err := ioutil.TempDir("", "peanut")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	err = source.fs.FetchAll(dir)
	if err != nil {
		return err
	}

	// Match files and copy to the destinations
	for _, mapping := range source.mappings {
		files, err := matchFiles(dir, mapping.MatchPattern)
		if err != nil {
			return err
		}

		fullDestPath := path.Join(destDir, mapping.DestPath)
		fsPath := filepath.FromSlash(fullDestPath)

		// If only a single file was matched specifically (not by glob matching),
		// treat the mapping.DestPath as the destination filename.
		if len(files) == 1 && !files[0].matchedByGlob {
			file := files[0]

			fileIsDir, err := file.IsDir()
			if err != nil {
				return err
			}

			if !fileIsDir {
				return file.CopyTo(fsPath)
			}
		}

		// Otherwise, copy all files to the destination directory
		for _, file := range files {
			err = file.CopyToDirectory(fsPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
