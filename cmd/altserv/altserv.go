package main

import (
	"net/http"

	"github.com/i2bskn/altserv"
)

func main() {
	config := altserv.NewConfig()
	handler := altserv.NewAppHandler(config)
	config.Logger.Printf("Server started with http://localhost%v", config.Port)
	http.ListenAndServe(config.Port, handler)
}
