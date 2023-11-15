package normalize

import "strings"

func NormalizePhoneNumber(phoneNumber string) string {
	var b strings.Builder
	for _, c := range phoneNumber {
		if '0' <= c && c <= '9' {
			b.WriteRune(c)
		}
	}
	return b.String()
}
