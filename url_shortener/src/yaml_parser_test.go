package urlshortener

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type testCaseYAML struct {
	data         string
	structFormat []URLMappingYAML
	mapFormat    map[string]string
}

var yamlTestCases = []testCaseYAML{
	{
		data: strings.Join([]string{
			"- path: /deimvis",
			"  url: https://github.com/Deimvis",
			"- path: /dbrusenin",
			"  url: https://www.linkedin.com/in/dmitriy-brusenin/",
		}, "\n"),
		structFormat: []URLMappingYAML{
			{Path: "/deimvis", URL: "https://github.com/Deimvis"},
			{Path: "/dbrusenin", URL: "https://www.linkedin.com/in/dmitriy-brusenin/"},
		},
		mapFormat: map[string]string{
			"/deimvis":   "https://github.com/Deimvis",
			"/dbrusenin": "https://www.linkedin.com/in/dmitriy-brusenin/",
		},
	},
}

func TestParseYAMLSimple(t *testing.T) {
	yamlFilePath := filepath.Join(t.TempDir(), "mapping.yaml")
	for _, testCase := range yamlTestCases {
		makeFile(t, yamlFilePath, testCase.data)
		result, err := parseYAML(yamlFilePath)
		if err != nil {
			t.Errorf("parseYAML failed: %s", err.Error())
		}
		if !reflect.DeepEqual(result, testCase.mapFormat) {
			t.Errorf("parseYAML returned unexpected result\nResult:   %s,\nExpected: %s\n", result, testCase.mapFormat)
		}
	}
}
func TestParseYAMLBadFilePath(t *testing.T) {
	_, err := parseYAML("nonexistentyamlfile.XblLqZdbiZlw4v2B")
	if err == nil {
		t.Errorf("parseYAML didn't return error on file that doesn't exist")
	}
}

func TestParseYAMLBadFileContent(t *testing.T) {
	yamlFilePath := filepath.Join(t.TempDir(), "mapping.yaml")
	makeFile(t, yamlFilePath, "bad yaml content")
	_, err := parseYAML(yamlFilePath)
	if err == nil {
		t.Errorf("parseYAML didn't return error on file with incorrect yaml content")
	}
}

func TestReadYAMLSimple(t *testing.T) {
	yamlFilePath := filepath.Join(t.TempDir(), "mapping.yaml")
	for _, testCase := range yamlTestCases {
		makeFile(t, yamlFilePath, testCase.data)
		result, err := readYAML(yamlFilePath)
		if err != nil {
			t.Errorf("readYAML failed: %s", err.Error())
		}
		if !reflect.DeepEqual(result, testCase.structFormat) {
			t.Errorf("readYAML returned unexpected result\nResult:   %s,\nExpected: %s\n", result, testCase.structFormat)
		}
	}
}

func TestReadYAMLBadFilepath(t *testing.T) {
	_, err := readYAML("nonexistentyamlfile.XblLqZdbiZlw4v2B")
	if err == nil {
		t.Errorf("readYAML didn't return error on file that doesn't exist")
	}
}

func TestReadYAMLBadFileContent(t *testing.T) {
	yamlFilePath := filepath.Join(t.TempDir(), "mapping.yaml")
	makeFile(t, yamlFilePath, "bad yaml content")
	_, err := readYAML(yamlFilePath)
	if err == nil {
		t.Errorf("readYAML didn't return error on file with incorrect yaml content")
	}
}
