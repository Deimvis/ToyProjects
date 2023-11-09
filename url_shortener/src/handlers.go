package urlshortener

import (
	"fmt"
	"net/http"
)

func NewMapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func NewYAMLHandler(filePath string, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := parseYAML(filePath)
	if err != nil {
		return nil, err
	}
	return NewMapHandler(pathsToUrls, fallback), nil
}

func NewJSONHandler(filePath string, fallback http.Handler) (http.HandlerFunc, error) {
	pathToUrls, err := parseJSON(filePath)
	if err != nil {
		return nil, err
	}
	return NewMapHandler(pathToUrls, fallback), nil
}

func NewPostgresHandler(tableName string, fallback http.Handler) (http.HandlerFunc, error) {
	if !isSafeTableName(tableName) {
		return nil, fmt.Errorf("got not safe table name: `%s` (please, do not inject anything)", tableName)
	}
	pathToUrls, err := parsePostgres(tableName)
	if err != nil {
		return nil, err
	}
	return NewMapHandler(pathToUrls, fallback), nil
}
