package main

import (
	"net/http"
	"os"
)

func initializeApp(config *Config) {
	if err := os.Mkdir(config.TmpDir, 0777); err != nil {
		config.Logger.Printf("TmpDir is already exist: %v", config.TmpDir)
	}
}

func main() {
	config := newConfig()
	initializeApp(config)
	handler := newAppHandler(config)
	config.Logger.Printf("Server started with http://localhost%v", config.Port)
	http.ListenAndServe(config.Port, handler)
}
