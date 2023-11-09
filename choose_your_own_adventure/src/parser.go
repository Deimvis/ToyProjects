package cyoa

import (
	"encoding/json"
	"io"
	"os"
)

func ParseJsonStoryFromFile(filePath string) (Story, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return Story{}, err
	}
	return ParseJsonStory(f)
}

func ParseJsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	err := d.Decode(&story)
	if err != nil {
		return Story{}, err
	}
	return story, nil
}
