package peanut

import (
	"github.com/bmatcuk/doublestar"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// localAwareFile remembers the baseDir it was matched on by matchFiles()
type localAwareFile struct {
	path     string
	basePath string
}

func newLocalAwareFile(path, basePath string) localAwareFile {
	return localAwareFile{path, basePath}
}

func (laFile *localAwareFile) RelativePath() (string, error) {
	return filepath.Rel(laFile.basePath, laFile.path)
}

func (laFile *localAwareFile) CopyToDirectory(dest string) error {
	fileInfo, err := fs.Stat(laFile.path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return nil
	}

	relativePath, err := laFile.RelativePath()
	if err != nil {
		return err
	}

	fullDestPath := filepath.Join(dest, relativePath)
	fullDestDir := filepath.Dir(fullDestPath)

	err = fs.MkdirAll(fullDestDir, os.ModePerm)
	if err != nil {
		return err
	}

	return copyFile(laFile.path, fullDestPath)
}

//func (laFile *localAwareFile) CopyTo(destPath string) error {
//	fileInfo, err := fs.Stat(laFile.FsPath())
//	if err != nil {
//		return err
//	}
//
//	fullDestPath := filepath.Join(destPath, laFile.relativePath)
//	fullDestDir := filepath.Dir(fullDestPath)
//
//	err = fs.MkdirAll(fullDestDir, os.ModePerm)
//	if err != nil {
//		return err
//	}
//
//	if fileInfo.IsDir() {
//		err = copyDirectory(laFile.FsPath(), fullDestPath)
//	} else {
//		err = copyFile(laFile.FsPath(), fullDestPath)
//	}
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

// matchFile returns a list of localAwareFile which were found in the baseDir based on the matchPattern.
// baseDir should be an os-specific path.
// The matchPattern is a forward-slash path representing one of the following:
// - glob pattern ("/usr/data/**/*.md")
// - file or directory path (/"usr/data/products" or "/usr/data/products/spark.md")
func matchFiles(baseDir, matchPattern string) ([]localAwareFile, error) {
	// If a wildcard character is present, we'll do glob matching
	if strings.Contains(matchPattern, "*") {
		return globMatchFiles(baseDir, matchPattern)
	}

	fullMatchPath := filepath.Join(baseDir, filepath.FromSlash(matchPattern))

	// Check if the file/directory exists
	fileInfo, err := fs.Stat(fullMatchPath)
	if err != nil {
		return nil, err
	}

	// If a single file matches, the basePath should be the directory of the file.
	if !fileInfo.IsDir() {
		return []localAwareFile{newLocalAwareFile(fullMatchPath, filepath.Dir(fullMatchPath))}, nil
	}

	return []localAwareFile{}, nil
}

func globMatchFiles(baseDir, matchPattern string) ([]localAwareFile, error) {
	var matchedFiles []localAwareFile

	globPath := filepath.Join(filepath.FromSlash(baseDir), filepath.FromSlash(matchPattern))
	globBase := getGlobBase(globPath)

	err := afero.Walk(fs, baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		match, err := doublestar.PathMatch(globPath, path)
		if err != nil {
			return err
		}

		if match {
			file := newLocalAwareFile(path, globBase)
			matchedFiles = append(matchedFiles, file)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matchedFiles, nil
}

// getGlobBase returns the path built from all the parts but ones containing glob elements.
// i.e: getGlobBase("/usr/data/**/*.md") = "/usr/data"
func getGlobBase(globPath string) string {
	var globBaseParts []string
	pathParts := strings.Split(globPath, string(os.PathSeparator))

	if strings.HasPrefix(globPath, "/") {
		globBaseParts = append(globBaseParts, "/")
	}

	for _, pathComponent := range pathParts {
		if strings.Contains(pathComponent, "*") {
			break
		}

		globBaseParts = append(globBaseParts, pathComponent)
	}

	return filepath.Join(globBaseParts...)
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
