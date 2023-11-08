package urlshortener

import (
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

func NewYAMLHandler(filepath string, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := parseYAML(filepath)
	if err != nil {
		return nil, err
	}
	return NewMapHandler(pathsToUrls, fallback), nil
}

func NewJSONHandler(filepath string, fallback http.Handler) (http.HandlerFunc, error) {
	pathToUrls, err := parseJSON(filepath)
	if err != nil {
		return nil, err
	}
	return NewMapHandler(pathToUrls, fallback), nil
}
