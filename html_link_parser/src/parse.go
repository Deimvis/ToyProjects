package hlp

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Takes reader of an HTML doc as argument,
// returns all links inside or an error
func ParseLinks(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	var links []Link
	dfs(doc, &links)
	return links, nil
}

func dfs(n *html.Node, links *[]Link) string {
	var textB strings.Builder
	if n.Type == html.TextNode {
		textB.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		innerText := dfs(c, links)
		if !isAtag(c) {
			textB.WriteString(innerText)
		}
	}
	text := textB.String()
	if isAtag(n) {
		link := newLink(n)
		link.Text = text
		*links = append(*links, *link)
	}
	return text
}

func isAtag(n *html.Node) bool {
	return n != nil && n.Type == html.ElementNode && n.Data == "a"
}

func newLink(n *html.Node) *Link {
	var l Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			l.Href = attr.Val
		}
	}
	return &l
}
