package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

const AppName = "AltServ"
const Index = "index.html"

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
		fmt.Println(err)
		return asset, err
	}

	if asset_info.IsDir() {
		asset = path.Join(asset, c.Index)
		asset_info, err = os.Stat(asset)
		if err != nil {
			fmt.Println(err)
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

func newConfig() *Config {
	return &Config{
		DocumentRoot: currentDir(),
		Logger:       generateLogger(AppName),
		Index:        Index,
	}
}
