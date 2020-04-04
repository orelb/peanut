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

	sources, err := peanut.CreateSources(config)
	if err != nil {
		log.Fatalf("Failed to create sources from configuration: %s", err)
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

	log.Println("Done :)")
}
