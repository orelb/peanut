package peanut

import (
	"gopkg.in/yaml.v2"
)

// SourceDeclaration defines a source to be used in Peanut
type SourceDeclaration struct {
	Name         string
	Type         string
	URL          string
	Revision     *string
	FileMappings []string `yaml:"files"`
}

// GetRevision returns the source declaration's revision if set or "master" if not.
func (sd SourceDeclaration) GetRevision() string {
	if sd.Revision == nil {
		return "master"
	}

	return *sd.Revision
}

// Config holds Peanut configuration
type Config struct {
	SourceDeclarations []SourceDeclaration `yaml:"sources"`
}

// ParseConfig loads Peanut configuration from a YAML string
func ParseConfig(configContents []byte) (*Config, error) {
	config := Config{}

	err := yaml.Unmarshal(configContents, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// CreateSources creates Peanut sources from Config.
func CreateSources(config *Config) ([]*Source, error) {
	sources := make([]*Source, len(config.SourceDeclarations))

	for i, sourceDeclaration := range config.SourceDeclarations {
		mappings := make([]FileMapping, len(sourceDeclaration.FileMappings))

		for j, fileMappingStr := range sourceDeclaration.FileMappings {
			parsedMapping, err := parseFileMapping(fileMappingStr)

			if err != nil {
				return nil, err
			}

			mappings[j] = parsedMapping
		}

		fs := NewGenericGitSourceFS(sourceDeclaration.URL, sourceDeclaration.GetRevision())
		source := NewSource(sourceDeclaration.Name, fs, mappings)

		sources[i] = source
	}

	return sources, nil
}
