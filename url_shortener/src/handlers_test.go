package urlshortener

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestNewMapHandlerSimple(t *testing.T) {
	mapHandler := NewMapHandler(map[string]string{"/path": "https://redirect.url/"}, http.NewServeMux())
	req, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatalf("Failed to create new request:\n%s\n", err.Error())
	}
	rr := httptest.NewRecorder()
	mapHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusFound {
		t.Errorf("MapHandler returned unexpected status code\nResult: %d\nExpected: %d\n", rr.Code, http.StatusFound)
	}
	if rr.Header().Get("Location") != "https://redirect.url/" {
		t.Errorf("MapHandler returned unexpected redirect url\nResult: %s\nExpected: %s\n", rr.Header().Get("Location"), "https://redirect.url/")
	}
}

func TestNewMapHandlerFallback(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/fallback_path", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "Hello, World!\n") })
	mapHandler := NewMapHandler(map[string]string{"/path": "https://redirect.url/"}, mux)
	req, err := http.NewRequest("GET", "/fallback_path", nil)
	if err != nil {
		t.Fatalf("Failed to create new request:\n%s\n", err.Error())
	}
	rr := httptest.NewRecorder()
	mapHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("MapHandler fallback returned unexpected status code\nResult: %d\nExpected: %d\n", rr.Code, http.StatusOK)
	}
	responseBody, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Failed to read response body:\n%s\n", err.Error())
	}
	if string(responseBody) != "Hello, World!\n" {
		t.Errorf("MapHandler fallback returned unexpected body\nResult: %q\nExpected: %q\n", string(responseBody), "Hello, World!\n")
	}
}

func TestNewYAMLHandlerSimple(t *testing.T) {
	yamlFilePath := filepath.Join(t.TempDir(), "mapping.yaml")
	for _, testCase := range yamlTestCases {
		makeFile(t, yamlFilePath, testCase.data)
		_, err := NewYAMLHandler(yamlFilePath, http.NewServeMux())
		if err != nil {
			t.Errorf("parseYAML failed: %s", err.Error())
		}
	}
}

func TestNewYAMLHandlerBadFilePath(t *testing.T) {
	_, err := NewYAMLHandler("nonexistentyamlfile.XblLqZdbiZlw4v2B", http.NewServeMux())
	if err == nil {
		t.Errorf("NewYAMLHandler didn't return error on file that doesn't exist")
	}
}

func TestNewYAMLHandlerBadFileContent(t *testing.T) {
	yamlFilePath := filepath.Join(t.TempDir(), "mapping.yaml")
	makeFile(t, yamlFilePath, "bad yaml content")
	_, err := NewYAMLHandler(yamlFilePath, http.NewServeMux())
	if err == nil {
		t.Errorf("NewYAMLHandler didn't return error on file with incorrect yaml content")
	}
}

func TestNewJSONHandlerSimple(t *testing.T) {
	jsonFilePath := filepath.Join(t.TempDir(), "mapping.json")
	for _, testCase := range jsonTestCases {
		makeFile(t, jsonFilePath, testCase.data)
		_, err := NewJSONHandler(jsonFilePath, http.NewServeMux())
		if err != nil {
			t.Errorf("parseJSON failed: %s", err.Error())
		}
	}
}

func TestNewJSONHandlerBadFilePath(t *testing.T) {
	_, err := NewJSONHandler("nonexistentjsonfile.XblLqZdbiZlw4v2B", http.NewServeMux())
	if err == nil {
		t.Errorf("NewJSONHandler didn't return error on file that doesn't exist")
	}
}

func TestNewJSONHandlerBadFileContent(t *testing.T) {
	jsonFilePath := filepath.Join(t.TempDir(), "mapping.json")
	makeFile(t, jsonFilePath, "bad json content")
	_, err := NewJSONHandler(jsonFilePath, http.NewServeMux())
	if err == nil {
		t.Errorf("NewJSONHandler didn't return error on file with incorrect json content")
	}
}
