package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	urlshortener "github.com/Deimvis/toyprojects/url_shortener/src"
)

func main() {
    jsonFilePath := flag.String("json", "", "Path to json url mapping file")
	yamlFilepath := flag.String("yaml", "", "Path to yaml url mapping file")
	flag.Parse()

	mux := initDefaultMux()
	handler := initDefaultHandler(mux)
	if *yamlFilepath != "" {
		yamlHandler, err := urlshortener.NewYAMLHandler(*yamlFilepath, handler)
		if err != nil {
			fmt.Printf("Failed to create YAML Handler:\n%s\n", err.Error())
			os.Exit(1)
		}
		handler = yamlHandler
	}
    if *jsonFilePath != "" {
        jsonHandler, err := urlshortener.NewJSONHandler(*jsonFilePath, handler)
        if err != nil {
            fmt.Printf("Failed to create JSON Handler:\n%s\n", err.Error())
            os.Exit(1)
        }
        handler = jsonHandler
    }
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func initDefaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Hello, World!") })
	return mux
}

func initDefaultHandler(fallback http.Handler) http.HandlerFunc {
	defaultPathsToUrls := map[string]string{
		"/google": "https://www.google.com/",
		"/github": "https://github.com/",
	}
	return urlshortener.NewMapHandler(defaultPathsToUrls, fallback)
}
