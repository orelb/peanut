package peanut

import (
	"errors"
	"fmt"
	"strings"
)

type FileMapping struct {
	MatchPattern string
	DestPath     string
}

func NewFileMapping(matchPattern, destPath string) FileMapping {
	return FileMapping{matchPattern, destPath}
}

// Parses "<sourceMatchPath>:<destinationPath>" type strings to a FileMapping instance.
func ParseFileMapping(fileMappingStr string) (FileMapping, error) {
	splitString := strings.Split(fileMappingStr, ":")

	if len(splitString) != 2 {
		return FileMapping{}, errors.New(fmt.Sprintf("Invalid file mapping string: '%s'", fileMappingStr))
	}

	matchPattern := splitString[0]
	destPath := splitString[1]

	return FileMapping{matchPattern, destPath}, nil
}
