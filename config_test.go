package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCurrentDir(t *testing.T) {
	expected_path, err := filepath.Abs(".")

	if err != nil {
		t.Fatal(err)
	}

	if currentDir() != expected_path {
		t.Fatalf("Expected %v, but %v", expected_path, currentDir())
	}
}

func TestDocumentRoot(t *testing.T) {
	if documentRoot() != currentDir() {
		t.Fatalf("Expected %v, but %v", currentDir(), documentRoot())
	}

	specific_document_root := "/path/to/docroot"
	os.Setenv(EnvDocRoot, specific_document_root)
	defer os.Setenv(EnvDocRoot, "")

	if documentRoot() != specific_document_root {
		t.Fatalf("Expected %v, but %v", specific_document_root, documentRoot())
	}
}

func TestNewConfig(t *testing.T) {
	config := newConfig()

	if config.DocumentRoot != documentRoot() {
		t.Fatalf("Expected %v, but %v", documentRoot(), config.DocumentRoot)
	}

	if config.Index != DefaultIndex {
		t.Fatalf("Expected %v, but %v", DefaultIndex, config.Index)
	}

	if config.TmpDir != DefaultTmpDir {
		t.Fatalf("Expected %v, but %v", DefaultTmpDir, config.TmpDir)
	}

	if config.Logger == nil {
		t.Fatalf("Expected %v, but %v", "*log.Logger", nil)
	}
}
