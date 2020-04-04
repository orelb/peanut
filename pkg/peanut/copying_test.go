package peanut

import (
	"bytes"
	"github.com/spf13/afero"
	"os"
	"testing"
)

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