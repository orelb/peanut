package peanut

import (
	"github.com/bmatcuk/doublestar"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
)

func MatchFiles(baseDir, matchPattern string) ([]LocalAwareFile, error) {
	// If a wildcard character is present, we'll do glob matching
	if strings.Contains(matchPattern, "*") {
		return globMatchFiles(baseDir, matchPattern)
	}

	fullMatchPath := filepath.Join(baseDir, filepath.FromSlash(matchPattern))

	// Check if the file/directory exists
	_, err := AppFs.Stat(fullMatchPath)
	if err != nil {
		return nil, err
	}

	// Convert the matchPattern to a relative pattern before initializing local-aware file
	relativePath, err := filepath.Rel(baseDir, filepath.Join(baseDir, filepath.FromSlash(matchPattern)))
	if err != nil {
		return nil, err
	}

	return []LocalAwareFile{newLocalAwareFile(baseDir, relativePath)}, nil
}

func globMatchFiles(baseDir, matchPattern string) ([]LocalAwareFile, error) {
	var matchedFiles []LocalAwareFile
	fullMatchPattern := filepath.Join(filepath.FromSlash(baseDir), filepath.FromSlash(matchPattern))

	err := afero.Walk(AppFs, baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		match, err := doublestar.PathMatch(fullMatchPattern, path)
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		if match {
			file := newLocalAwareFile(baseDir, relativePath)
			matchedFiles = append(matchedFiles, file)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matchedFiles, nil
}
