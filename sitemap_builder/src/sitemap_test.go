package sitemap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeBaseUrl(t *testing.T) {
	testCases := []struct {
		url      string
		expected string
	}{
		{"https://bebest.pro/posts/?speciality=backend", "https://bebest.pro"},
		{"https://bebest.pro/posts/?speciality=backend&amp;box_animation=true", "https://bebest.pro"},
		{"https://github.com/Deimvis", "https://github.com"},
	}
	for _, tc := range testCases {
		baseUrl, err := makeBaseUrl(tc.url)
		require.NoError(t, err)
		require.Equal(t, tc.expected, baseUrl)
	}
}
