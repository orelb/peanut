package peanut

import "gopkg.in/yaml.v2"

type SourceDeclaration struct {
	Name          string
	Type          string
	RepositoryURL string `yaml:"repository_url"`
	FileMappings  map[string]string `yaml:"files"`
}

type Config struct {
	SourceDeclarations []SourceDeclaration `yaml:"sources"`
}

// Loads Peanut configuration from YAML string
func ParseConfig(configContents []byte) (*Config, error) {
	config := Config{}

	err := yaml.Unmarshal(configContents, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
