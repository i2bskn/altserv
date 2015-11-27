package altserv

import (
	"log"
	"os"
	"path/filepath"
)

const (
	AppName      = "AltServ"
	EnvDocRoot   = "AS_DROOT"
	DefaultPort  = ":10080"
	DefaultIndex = "index.html"
)

type Config struct {
	DocumentRoot string
	Port         string
	Index        string
	Logger       *log.Logger
}

func currentDir() string {
	currentPath, err := filepath.Abs(".")
	if err != nil {
		panic("Current path can not obtain.")
	}
	return currentPath
}

func documentRoot() string {
	path := os.Getenv(EnvDocRoot)
	if len(path) == 0 {
		return currentDir()
	}
	return path
}

func NewConfig() *Config {
	return &Config{
		DocumentRoot: documentRoot(),
		Port:         DefaultPort,
		Index:        DefaultIndex,
		Logger:       generateLogger(AppName),
	}
}
