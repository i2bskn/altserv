package main

import (
	"log"
	"os"
	"path/filepath"
)

const AppName = "AltServ"
const EnvDocRoot = "AS_DROOT"
const DefaultIndex = "index.html"
const DefaultTmpDir = "/tmp/altserv_temporary"

type Config struct {
	DocumentRoot string
	Index        string
	TmpDir       string
	Logger       *log.Logger
}

func currentDir() string {
	current_path, err := filepath.Abs(".")
	if err != nil {
		panic("Current path can not obtain.")
	}
	return current_path
}

func documentRoot() string {
	path := os.Getenv(EnvDocRoot)
	if len(path) == 0 {
		return currentDir()
	}
	return path
}

func newConfig() *Config {
	return &Config{
		DocumentRoot: documentRoot(),
		Index:        DefaultIndex,
		TmpDir:       DefaultTmpDir,
		Logger:       generateLogger(AppName),
	}
}
