package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

const AppName = "AltServ"
const Index = "index.html"
const EnvDocRoot = "AS_DROOT"

type Config struct {
	DocumentRoot string
	Index        string
	Logger       *log.Logger
}

func (c *Config) AssetPath(uri string) (asset string, err error) {
	var asset_info os.FileInfo
	asset = path.Join(c.DocumentRoot, uri)
	asset_info, err = os.Stat(asset)

	if err != nil {
		return asset, err
	}

	if asset_info.IsDir() {
		asset = path.Join(asset, c.Index)
		asset_info, err = os.Stat(asset)
		if err != nil {
			return asset, err
		}
	}

	return asset, nil
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
		fmt.Println("docroot from current dir")
		return currentDir()
	}
	return path
}

func newConfig() *Config {
	return &Config{
		DocumentRoot: documentRoot(),
		Logger:       generateLogger(AppName),
		Index:        Index,
	}
}
