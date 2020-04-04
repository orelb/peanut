package peanut

import (
	"fmt"
	"strings"
)

// FileMapping declares a mapping between matched files (by the MatchPattern) to their destination path.
type FileMapping struct {
	MatchPattern string
	DestPath     string
}

// NewFileMapping creates a new FileMapping instance
func NewFileMapping(matchPattern, destPath string) FileMapping {
	return FileMapping{matchPattern, destPath}
}

// parseFileMapping parses "<sourceMatchPath>:<destinationPath>" type strings to a FileMapping instance.
func parseFileMapping(fileMappingStr string) (FileMapping, error) {
	splitString := strings.Split(fileMappingStr, ":")

	if len(splitString) != 2 {
		return FileMapping{}, fmt.Errorf("invalid file mapping string: '%s'", fileMappingStr)
	}

	matchPattern := splitString[0]
	destPath := splitString[1]

	return FileMapping{matchPattern, destPath}, nil
}
