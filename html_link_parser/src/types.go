package hlp

import (
	"fmt"
	"strings"
)

// Link represents <a href="..."> in an HTML
type Link struct {
	Href string
	Text string
}

func (l Link) ToString() string {
	return fmt.Sprintf("%s: %q", l.Href, strings.TrimSpace(l.Text))
}
