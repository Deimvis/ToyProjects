package normalize

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizePhoneNumber(t *testing.T) {
	testCases := []struct {
		phoneNumber string
		expected    string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"123-456-7890", "1234567890"},
		{"1234567892", "1234567892"},
		{"(123)456-7892", "1234567892"},
	}
	for _, tc := range testCases {
		t.Run(tc.phoneNumber, func(t *testing.T) {
			require.Equal(t, tc.expected, NormalizePhoneNumber(tc.phoneNumber))
		})
	}
}
