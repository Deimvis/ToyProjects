package main

import (
	"flag"
	"fmt"
	"os"

	hlp "github.com/Deimvis/toyprojects/html_link_parser/src"
)

func main() {
	filePath := flag.String("file", "", "Path to an html file")
	flag.Parse()
	f, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n%s\n", *filePath, err.Error())
		os.Exit(1)
	}
	defer f.Close()
	links, err := hlp.ParseLinks(f)
	if err != nil {
		fmt.Printf("Failed to parse links\n%s\n", err.Error())
		os.Exit(1)
	}
	for _, l := range links {
		fmt.Println(l.ToString())
	}
}
