package urlshortener

import (
	"path/filepath"
	"reflect"
	"testing"
)

type testCaseJSON struct {
	data      string
	mapFormat map[string]string
}

var jsonTestCases = []testCaseJSON{
	{
		data: `{
			"/deimvis": "https://github.com/Deimvis",
			"/dbrusenin": "https://www.linkedin.com/in/dmitriy-brusenin/"
		}`,
		mapFormat: map[string]string{
			"/deimvis":   "https://github.com/Deimvis",
			"/dbrusenin": "https://www.linkedin.com/in/dmitriy-brusenin/",
		},
	},
}

func TestParseJSONSimple(t *testing.T) {
	jsonFilePath := filepath.Join(t.TempDir(), "mapping.json")
	for _, testCase := range jsonTestCases {
		makeFile(t, jsonFilePath, testCase.data)
		result, err := parseJSON(jsonFilePath)
		if err != nil {
			t.Errorf("parseJSON failed: %s", err.Error())
		}
		if !reflect.DeepEqual(result, testCase.mapFormat) {
			t.Errorf("parseJSON returned unexpected result\nResult:   %s,\nExpected: %s\n", result, testCase.mapFormat)
		}
	}
}
func TestParseJSONBadFilePath(t *testing.T) {
	_, err := parseJSON("nonexistentJSONfile.XblLqZdbiZlw4v2B")
	if err == nil {
		t.Errorf("parseJSON didn't return error on file that doesn't exist")
	}
}

func TestParseJSONBadFileContent(t *testing.T) {
	jsonFilePath := filepath.Join(t.TempDir(), "mapping.json")
	makeFile(t, jsonFilePath, "bad JSON content")
	_, err := parseJSON(jsonFilePath)
	if err == nil {
		t.Errorf("parseJSON didn't return error on file with incorrect JSON content")
	}
}
