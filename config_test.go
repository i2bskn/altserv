package altserv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCurrentDir(t *testing.T) {
	expected, err := filepath.Abs(".")

	if err != nil {
		t.Fatal(err)
	}

	if currentDir() != expected {
		t.Fatalf("Expected %v, but %v", expected, currentDir())
	}
}

func TestDocumentRoot(t *testing.T) {
	if documentRoot() != currentDir() {
		t.Fatalf("Expected %v, but %v", currentDir(), documentRoot())
	}

	expected := "/path/to/docroot"
	os.Setenv(EnvDocRoot, expected)
	defer os.Unsetenv(EnvDocRoot)

	if documentRoot() != expected {
		t.Fatalf("Expected %v, but %v", expected, documentRoot())
	}
}

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.DocumentRoot != documentRoot() {
		t.Fatalf("Expected %v, but %v", documentRoot(), config.DocumentRoot)
	}

	if config.Index != DefaultIndex {
		t.Fatalf("Expected %v, but %v", DefaultIndex, config.Index)
	}

	if config.Logger == nil {
		t.Fatalf("Expected %v, but %v", "*log.Logger", nil)
	}
}
