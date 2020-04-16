package main

import (
	"github.com/orelb/peanut/pkg/peanut"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var configFile string

	app := &cli.App{
		Name:  "Peanut",
		Usage: "fetch files from different sources to help compose your static website.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-file",
				Aliases:     []string{"c"},
				Value:       "peanut.yml",
				Usage:       "path to peanut config file",
				Destination: &configFile,
			},
		},
		Action: func(context *cli.Context) error {
			pullAllSources(configFile)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func pullAllSources(configFile string) {
	configContents, err := ioutil.ReadFile(configFile)
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
