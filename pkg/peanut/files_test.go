package peanut

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/afero"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func CreateFiles(filenames []string) {
	fileContents := []byte("some text")

	for _, filename := range filenames {
		fs.MkdirAll(path.Dir(filename), os.ModePerm)
		afero.WriteFile(fs, filepath.FromSlash(filename), fileContents, os.ModePerm)
	}
}

func InitTestFiles() {
	fs = afero.NewMemMapFs()

	CreateFiles([]string{
		"/data/a.md",
		"/data/b.md",
		"/data/recipes/readme.md",
		"/data/recipes/italian/pasta.md",
		"/data/recipes/italian/pasta-tasty.png",
		"/data/recipes/italian/pizza.md",
		"/data/recipes/italian/pizza1.png",
		"/data/recipes/italian/pizza2.png",
		"/data/recipes/israeli/falafel.md",
		"/data/recipes/israeli/falafel-with-tahini.png",
	})
}

func createLocalAwareFileFromPaths(baseDir string, paths []string) []localAwareFile {
	files := make([]localAwareFile, len(paths))

	for i, path := range paths {
		files[i] = newLocalAwareFile(baseDir, path)
	}

	return files
}

var localAwareFileComparer = cmp.Comparer(func(x, y localAwareFile) bool {
	return x.path == y.path && x.basePath == y.basePath
})

var sortLocalAwareFiles = cmpopts.SortSlices(func(x, y localAwareFile) bool { return x.path < y.path })

type MatchFileTest struct {
	matchPattern  string
	baseDir       string
	expectedPaths []string
}

var matchFilesTests = []MatchFileTest{

	// Glob tests
	{"/data/recipes/**/*.png", "/",
		[]string{
			"data/recipes/italian/pasta-tasty.png",
			"data/recipes/italian/pizza1.png",
			"data/recipes/italian/pizza2.png",
			"data/recipes/israeli/falafel-with-tahini.png",
		},
	},
	{"/*", "/data", []string{"a.md", "b.md", "recipes"}},
	{"/**/*.md", "/data/recipes", []string{
		"readme.md", "israeli/falafel.md", "italian/pizza.md", "italian/pasta.md",
	}},

	// Files/Directories tests
	{"italian", "/data/recipes", []string{"italian"}},
	{"/data/a.md", "/", []string{"data/a.md"}},
	{"recipes", "/data", []string{"recipes"}},
}

func TestMatchFiles(t *testing.T) {
	InitTestFiles()

	for _, test := range matchFilesTests {
		testName := fmt.Sprintf("('%s','%s')", test.matchPattern, test.baseDir)
		t.Run(testName, func(t *testing.T) {
			expectedFiles := createLocalAwareFileFromPaths(test.baseDir, test.expectedPaths)
			matchedFiles, err := matchFiles(test.baseDir, test.matchPattern)
			if err != nil {
				t.Fatalf("Failed to get files: %s", err)
			}

			if diff := cmp.Diff(expectedFiles, matchedFiles, localAwareFileComparer, sortLocalAwareFiles); diff != "" {
				t.Fatalf("Matched paths mismatch. (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	fs = afero.NewMemMapFs()
	copiedFilePath := "/adsad"
	expectedData := []byte("This is some test data")

	_ = afero.WriteFile(fs, "/data", expectedData, os.ModePerm)

	err := copyFile("/data", copiedFilePath)
	if err != nil {
		t.Errorf("copyFile() failed: %s", err)
	}

	_, err = fs.Stat(copiedFilePath)
	if err != nil {
		t.Errorf("Failed to stat %s: %s", copiedFilePath, err)
	}

	copiedFileData, err := afero.ReadFile(fs, copiedFilePath)
	if bytes.Compare(copiedFileData, expectedData) != 0 {
		t.Errorf("Copied file data and original data is different")
	}
}
