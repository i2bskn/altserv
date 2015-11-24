package main

import (
	"log"
	"testing"
)

func TestGenerateLogger(t *testing.T) {
	app_name := "TEST"
	expected_prefix := "[" + app_name + "] "
	logger := generateLogger(app_name)

	if logger.Prefix() != expected_prefix {
		t.Fatalf("Expected %v, but %v", expected_prefix, logger.Prefix())
	}

	if logger.Flags() != log.LstdFlags {
		t.Fatalf("Expected %v, but %v", log.LstdFlags, logger.Flags())
	}
}
