package main

import (
	"log"
	"os"
)

func generateLogger(name string) *log.Logger {
	prefix := "[" + name + "] "
	return log.New(os.Stderr, prefix, log.LstdFlags)
}
