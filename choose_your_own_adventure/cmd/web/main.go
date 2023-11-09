package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	cyoa "github.com/Deimvis/toyprojects/choose_your_own_adventure/src"
)

func main() {
	port := flag.Int("port", 3000, "The port for web server")
	filePath := flag.String("file", "story.json", "Path to JSON file with Choose Your Story Adventure")
	flag.Parse()
	story, err := cyoa.ParseJsonStoryFromFile(*filePath)
	if err != nil {
		fmt.Printf("Parsing story from given JSON file failed:\n%s\n", err.Error())
		os.Exit(1)
	}
	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		fmt.Printf("Failed to make template of `static/index.html`:\n%s\n", err.Error())
		os.Exit(1)
	}
	h := cyoa.NewHandler(story, cyoa.WithTemplate(tmpl))
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
