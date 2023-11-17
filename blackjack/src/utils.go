package bj

import (
	"fmt"
	"strings"
)

func ljust(s string, n int, fill string) string {
	if len(fill) != 1 {
		panic(fmt.Errorf("ljust needs fill argument to have length 1, received: %d (%q)", len(fill), fill))
	}
	return s + strings.Repeat(fill, max(n-len(s), 0))
}
