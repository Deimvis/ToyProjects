package cyoa

import (
	"os"
	"testing"
)

func makeFile(t *testing.T, filePath string, content string) {
	t.Helper()
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Couldn't write to file `%s`: %s", filePath, err.Error())
	}
}
