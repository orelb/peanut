package peanut

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseFileMapping_Basic(t *testing.T) {
	expectedMapping := FileMapping{"/docs", "/product1/docs"}
	fileMappingStr := "/docs:/product1/docs"

	parsedMapping, err := parseFileMapping(fileMappingStr)
	if err != nil {
		t.Fatalf("Failed parsing file mapping: %s", err)
	}

	if !cmp.Equal(parsedMapping, expectedMapping) {
		t.Fatalf("Got %s, expected %s", parsedMapping, expectedMapping)
	}
}

func TestParseFileMapping_ExpectError(t *testing.T) {
	invalidInputs := []string {
		"/docs/product1/docs",
		"hello:to:here",
		"",
	}

	for _, invalidInput := range invalidInputs {
		t.Run(invalidInput, func(t *testing.T) {
			_, err := parseFileMapping(invalidInput)

			if err == nil {
				t.Fatalf("parseFileMapping should return error")
			}
		})
	}

}
