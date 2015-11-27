package main

import (
	"log"
	"testing"
)

func TestGenerateLogger(t *testing.T) {
	appName := "TEST"
	logger := generateLogger(appName)

	{
		expected := "[" + appName + "] "
		if logger.Prefix() != expected {
			t.Fatalf("Expected %v, but %v", expected, logger.Prefix())
		}
	}

	{
		expected := log.LstdFlags
		if logger.Flags() != expected {
			t.Fatalf("Expected %v, but %v", expected, logger.Flags())
		}
	}
}
