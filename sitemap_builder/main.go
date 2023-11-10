package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	sitemap "github.com/Deimvis/toyprojects/sitemap_builder/src"
)

func encodeXML(obj any) []byte {
	xmlBytes, err := xml.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Printf("Failed to encode sitemap into XML format:\n%s\n", err.Error())
		os.Exit(1)
	}
	xmlBytes = append([]byte(xml.Header), xmlBytes...)
	return xmlBytes
}

func main() {
	urlFlag := flag.String("url", "https://bebest.pro", "Url to build a sitemap for")
	maxDepth := flag.Int("depth", 3, "Max depth for page traversal")
	flag.Parse()

	sitemap, err := sitemap.BuildSitemap(*urlFlag, *maxDepth)
	if err != nil {
		fmt.Printf("Failed to build sitemap:\n%s\n", err.Error())
		os.Exit(1)
	}
	xmlBytes := encodeXML(sitemap)
	fmt.Println(string(xmlBytes))
}
