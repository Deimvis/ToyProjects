package urlshortener

import (
	"os"

	"gopkg.in/yaml.v3"
)

type URLMappingYAML struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func parseYAML(filepath string) (map[string]string, error) {
	fileContent, err := readYAML(filepath)
	if err != nil {
		return nil, err
	}
	pathToUrls := make(map[string]string)
	for _, mapping := range fileContent {
		pathToUrls[mapping.Path] = mapping.URL
	}
	return pathToUrls, nil
}

func readYAML(filepath string) ([]URLMappingYAML, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	var fileContent []URLMappingYAML
	err = decoder.Decode(&fileContent)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}
