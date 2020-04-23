package main

import (
	"fmt"
	"github.com/orelb/peanut/pkg/peanut"
	"github.com/urfave/cli/v2"
	"io/ioutil"
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
			return pullAllSources(configFile)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func pullAllSources(configFile string) error {
	configContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config file.\n%s", err)
	}

	config, err := peanut.ParseConfig(configContents)
	if err != nil {
		return fmt.Errorf("failed to parse config file.\n%s", err)
	}

	sources, err := peanut.CreateSources(config)
	if err != nil {
		return fmt.Errorf("failed to initialize sources from configuration.\n%s", err)
	}

	fmt.Printf("Created %d sources.\n", len(sources))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to fetch current working directory.\n%s", err)
	}

	for _, source := range sources {
		err := source.Fetch(cwd)
		fmt.Printf("Fetching source \"%s\"\n", source.Name())

		if err != nil {
			return fmt.Errorf("failed to pull source %s.\n%s", source.Name(), err)
		}
	}

	fmt.Println("Done.")
	return nil
}
