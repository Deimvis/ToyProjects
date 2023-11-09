package urlshortener

import "testing"

func TestIsSafeTableName(t *testing.T) {
	testCases := []struct {
		tableName string
		expected  bool
	}{
		{"user", true},
		{"url_mapping", true},
		{"DROP TABLE", false},
		{"DROP DATABASE", false},
	}
	for _, tCase := range testCases {
		result := isSafeTableName(tCase.tableName)
		if result != tCase.expected {
			t.Errorf("isSafeTableName returned unexpected result\nResult: %t\nExpected: %t\n", result, tCase.expected)
		}
	}
}
