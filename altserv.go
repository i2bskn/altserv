package main

import (
	"net/http"
)

const AppName = "AltServ"

func main() {
	handler := newAppHandler(AppName)
	handler.logger.Println("Server started with http://localhost:10080")
	http.ListenAndServe(":10080", handler)
}
