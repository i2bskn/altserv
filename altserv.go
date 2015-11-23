package main

import (
	"net/http"
	"os"
)

func initializeApp(config *Config) {
	os.Mkdir(config.TmpDir, 0777)
}

func main() {
	config := newConfig()
	initializeApp(config)
	handler := newAppHandler(config)
	config.Logger.Println("Server started with http://localhost:10080")
	http.ListenAndServe(":10080", handler)
}
