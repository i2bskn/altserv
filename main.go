package main

import "net/http"

func main() {
	config := newConfig()
	handler := newAppHandler(config)
	config.Logger.Printf("Server started with http://localhost%v", config.Port)
	http.ListenAndServe(config.Port, handler)
}
