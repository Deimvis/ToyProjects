package sitemap

import (
	"net/http"
	"net/url"
	"strings"

	hlp "github.com/Deimvis/toyprojects/html_link_parser/src"
)

const xmlns = "https://www.sitemaps.org/schemas/sitemap/0.9/"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Xmlns string `xml:"xmlns,attr"`
	Urls  []loc  `xml:"url"`
}

type Sitemap urlset

func BuildSitemap(urlStr string, maxDepth int) (Sitemap, error) {
	pages, err := gatherPages(urlStr, maxDepth)
	sitemap := Sitemap{
		Xmlns: xmlns,
		Urls:  make([]loc, len(pages)),
	}
	for i := range pages {
		sitemap.Urls[i] = loc{pages[i]}
	}
	return sitemap, err
}

// Returns unique pages having the same domain
func gatherPages(urlStr string, maxDepth int) (pages []string, err error) {
	seen := make(map[string]struct{})
	var dfs func(page string, depth int)
	dfs = func(page string, depth int) {
		seen[page] = struct{}{}
		if depth >= maxDepth-1 {
			return
		}
		otherPages, err := iteratePages(page)
		if err != nil {
			panic(err)
		}
		newPages := *filterInPlace(&otherPages, func(p string) bool {
			if _, ok := seen[p]; ok {
				return false
			}
			seen[p] = struct{}{}
			return true
		})
		for _, p := range newPages {
			dfs(p, depth+1)
		}
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	dfs(urlStr, 0)
	for p := range seen {
		pages = append(pages, p)
	}
	return
}

// Returns all pages found on `urlStr` having the same domain (possibly with duplicates)
func iteratePages(urlStr string) ([]string, error) {
	baseUrl, err := makeBaseUrl(urlStr)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	links, err := hlp.ParseLinks(resp.Body)
	if err != nil {
		return nil, err
	}
	var pages []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			pages = append(pages, baseUrl+l.Href)
		case strings.HasPrefix(l.Href, baseUrl):
			pages = append(pages, l.Href)
		}
	}
	return pages, nil
}

func makeBaseUrl(urlStr string) (string, error) {
	urlParsed, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	baseUrl := (&url.URL{
		Scheme: urlParsed.Scheme,
		Host:   urlParsed.Host,
	}).String()
	return baseUrl, nil
}
