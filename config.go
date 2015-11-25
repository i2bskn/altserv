package main

import (
	"log"
	"os"
	"path/filepath"
)

const (
	AppName       = "AltServ"
	EnvDocRoot    = "AS_DROOT"
	DefaultPort   = ":10080"
	DefaultIndex  = "index.html"
	DefaultTmpDir = "/tmp/altserv_temporary"
)

type Config struct {
	DocumentRoot string
	Port         string
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
		Port:         DefaultPort,
		Index:        DefaultIndex,
		TmpDir:       DefaultTmpDir,
		Logger:       generateLogger(AppName),
	}
}
