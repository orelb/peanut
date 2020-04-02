package main

import (
	"github.com/orelb/peanut/pkg/peanut"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	configContents, err := ioutil.ReadFile("peanut.yml")
	if err != nil {
		log.Fatalf("Failed to load config file: %s", err)
	}

	config, err := peanut.ParseConfig(configContents)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err)
	}

	sources := make([]*peanut.Source, len(config.SourceDeclarations))
	for i, sourceDeclaration := range config.SourceDeclarations {
		mappings := make([]peanut.FileMapping, len(sourceDeclaration.FileMappings))
		j := 0

		for sourceMatcher, destPath := range sourceDeclaration.FileMappings {
			mappings[j] = peanut.NewFileMapping(peanut.NewStandardFileSelector(sourceMatcher), destPath)
			j++
		}

		fs := peanut.NewGenericGitSourceFS(sourceDeclaration.RepositoryURL)
		source := peanut.NewSource(fs, mappings)

		sources[i] = source
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to fetch current working directory: %s", err)
	}

	for _, source := range sources {
		err := source.Pull(cwd)

		if err != nil {
			log.Fatalf("Failed to pull source: %s", err)
		}
	}
}
