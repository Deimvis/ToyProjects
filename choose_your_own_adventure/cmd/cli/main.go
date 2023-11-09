package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	cyoa "github.com/Deimvis/toyprojects/choose_your_own_adventure/src"
)

func main() {
	filePath := flag.String("file", "story.json", "Path to JSON file with Choose Your Story Adventure")
	flag.Parse()
	story, err := cyoa.ParseJsonStoryFromFile(*filePath)
	if err != nil {
		fmt.Printf("Parsing story from given JSON file failed:\n%s\n", err.Error())
		os.Exit(1)
	}
	runner := cyoa.NewCLIRunner(story)
	log.Fatal(runner.Run())
}
