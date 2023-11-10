package sitemap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterInPlaceSimple(t *testing.T) {
	testCases := []struct {
		arr      []string
		f        filterFn
		expected []string
	}{
		{
			[]string{"1", "22", "333", "4444", "55555"},
			func(s string) bool {
				return len(s)%2 == 0
			},
			[]string{"22", "4444"},
		},
	}
	for _, tc := range testCases {
		filterInPlace(&tc.arr, tc.f)
		require.Equal(t, tc.expected, tc.arr)
	}
}
