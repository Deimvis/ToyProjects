package urlshortener

import (
	"encoding/json"
	"os"
)

type URLMappingJSON struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func parseJSON(filepath string) (map[string]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var fileContent map[string]string
	err = decoder.Decode(&fileContent)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}
