package hlp

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestIsAtag(t *testing.T) {
	testCases := []struct {
		html     string
		expected bool
	}{
		{`<a href="">Simple Case</a>`, true},
		{`<div>Simple Case</div>`, false},
		{`<a>no href</a>`, true},
	}
	for _, tc := range testCases {
		node := buildNode(t, tc.html)
		require.Equal(t, tc.expected, isAtag(node), "html: %s", tc.html)
	}
}

func TestDfs(t *testing.T) {
	testCases := []struct {
		htmlFilePath string
		expected     []Link
	}{
		{"../test_data/simple.html", []Link{{"/other-page", "A link to another page"}}},
		{"../test_data/inner_tag.html", []Link{
			{"https://www.twitter.com/joncalhoun", "\n      Check me out on twitter\n      \n    "},
			{"https://github.com/gophercises", "\n      Gophercises is on Github!\n    "},
		}},
		{"../test_data/real_life.html", []Link{
			{"#", "Login "},
			{"/lost", "Lost? Need help?"},
			{"https://twitter.com/marcusolsson", "@marcusolsson"},
		}},
		{"../test_data/comments.html", []Link{{"/dog-cat", "dog cat "}}},
		{"../test_data/no_href.html", []Link{}},
	}
	for _, tc := range testCases {
		links := []Link{}
		dfs(buildDoc(t, tc.htmlFilePath), &links)
		require.Equal(t, tc.expected, links, "html file: %s", tc.htmlFilePath)
	}
}

func buildNode(t *testing.T, htmlStr string) *html.Node {
	t.Helper()
	reader := strings.NewReader(htmlStr)
	doc, err := html.Parse(reader)
	if err != nil {
		t.Fatalf("Failed to parse given htmlStr:\n%s\nGot error:%s\n", htmlStr, err.Error())
	}
	return doc.FirstChild.LastChild.FirstChild
}

func buildDoc(t *testing.T, htmlFilePath string) *html.Node {
	t.Helper()
	file, err := os.Open(htmlFilePath)
	if err != nil {
		t.Fatalf("Failed to open given file: %s\n%s\n", htmlFilePath, err.Error())
	}
	doc, err := html.Parse(file)
	if err != nil {
		t.Fatalf("Failed to parse given html file: %s\n%s\n", htmlFilePath, err.Error())
	}
	return doc
}

// func dfslog(t *testing.T, n *html.Node, depth int) {
// 	t.Log(strings.Repeat(" ", depth) + n.Data)
// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
// 		dfslog(t, c, depth+1)
// 	}
// }
