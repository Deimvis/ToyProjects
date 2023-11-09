package cyoa

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCaseJSON struct {
	data  string
	story Story
}

var jsonTestCases = []testCaseJSON{
	{
		data: `{
			"meta": {
				"first_chapter_title": "intro"
			},
			"chapters": {
				"intro": {
					"title": "Title",
					"story": ["Paragraph 1", "Paragraph 2", "Paragraph 3"],
					"options": [
						{"text": "Option1", "arc": "arc1"},
						{"text": "Option2", "arc": "arc2"}
					]
				}
			}
		}`,
		story: Story{
			Meta: Meta{FirstChapterTitle: "intro"},
			Chapters: map[string]Chapter{
				"intro": {
					Title:      "Title",
					Paragraphs: []string{"Paragraph 1", "Paragraph 2", "Paragraph 3"},
					Options: []Option{
						{"Option1", "arc1"},
						{"Option2", "arc2"},
					},
				},
			},
		},
	},
}

func TestParseJsonStoryFromFileSimple(t *testing.T) {
	jsonFilePath := filepath.Join(t.TempDir(), "story.json")
	for _, tc := range jsonTestCases {
		makeFile(t, jsonFilePath, tc.data)
		result, err := ParseJsonStoryFromFile(jsonFilePath)
		require.NoError(t, err)
		require.Equal(t, result, tc.story)
	}
}

func TestParseJSONBadFilePath(t *testing.T) {
	_, err := ParseJsonStoryFromFile("nonexistentJSONfile.XblLqZdbiZlw4v2B")
	require.Error(t, err)
}

func TestParseJSONBadFileContent(t *testing.T) {
	jsonFilePath := filepath.Join(t.TempDir(), "story.json")
	makeFile(t, jsonFilePath, "bad JSON content")
	_, err := ParseJsonStoryFromFile(jsonFilePath)
	require.Error(t, err)
}

func TestParseJsonStory(t *testing.T) {
	for _, tc := range jsonTestCases {
		r := strings.NewReader(tc.data)
		result, err := ParseJsonStory(r)
		require.NoError(t, err)
		require.Equal(t, result, tc.story)
	}
}

func testParseJsonStoryBadFileContent(t *testing.T) {
	r := strings.NewReader("bad JSON content")
	_, err := ParseJsonStory(r)
	require.Error(t, err)
}
