package peanut

import "gopkg.in/yaml.v2"

type SourceDeclaration struct {
	Name          string
	Type          string
	RepositoryURL string   `yaml:"repository_url"`
	FileMappings  []string `yaml:"files"`
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

func CreateSources(config *Config) ([]*Source, error) {
	sources := make([]*Source, len(config.SourceDeclarations))

	for i, sourceDeclaration := range config.SourceDeclarations {
		mappings := make([]FileMapping, len(sourceDeclaration.FileMappings))

		for j, fileMappingStr := range sourceDeclaration.FileMappings {
			parsedMapping, err := ParseFileMapping(fileMappingStr)

			if err != nil {
				return nil, err
			}

			mappings[j] = parsedMapping
		}

		fs := NewGenericGitSourceFS(sourceDeclaration.RepositoryURL)
		source := NewSource(fs, mappings)

		sources[i] = source
	}

	return sources, nil
}
