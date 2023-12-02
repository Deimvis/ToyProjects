package quiet_hn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCacheWithTTL(t *testing.T) {
	require.NotPanics(t, func() { NewCacheWithTTL[string, string]() })
}

func TestCacheWithTTL_PutGet(t *testing.T) {
	type put struct {
		key   string
		value string
	}
	type getResult struct {
		value string
		ok    bool
	}
	type get struct {
		key      string
		expected getResult
	}
	testCases := []struct {
		puts []put
		gets []get
	}{
		{
			puts: []put{{"key", "value"}},
			gets: []get{
				{"key", getResult{"value", true}},
				{"nokey", getResult{"???", false}},
				{"randomkey", getResult{"???", false}},
			},
		},
		{
			puts: []put{{"1", "1"}, {"2", "2"}},
			gets: []get{
				{"1", getResult{"1", true}},
				{"2", getResult{"2", true}},
				{"3", getResult{"???", false}},
			},
		},
	}
	for _, tc := range testCases {
		cache := NewCacheWithTTL[string, string]()
		for _, p := range tc.puts {
			cache.Put(p.key, p.value)
		}
		for _, g := range tc.gets {
			value, ok := cache.Get(g.key)
			require.Equal(t, g.expected.ok, ok)
			if ok {
				require.Equal(t, g.expected.value, value)
			}

		}
	}
}
