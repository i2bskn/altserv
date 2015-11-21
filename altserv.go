package main

import (
	"net/http"
)

func main() {
	config := newConfig()
	handler := newAppHandler(config)
	config.Logger.Println("Server started with http://localhost:10080")
	http.ListenAndServe(":10080", handler)
}
