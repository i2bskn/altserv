package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type AppHandler struct {
	logger *log.Logger
}

func (h *AppHandler) RequestLog(r *http.Request) {
	log_info := []string{
		r.Method,
		r.URL.Path,
	}

	h.logger.Println(strings.Join(log_info, " "))
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.RequestLog(r)
	io.WriteString(w, "example")
}

func newAppHandler(name string) *AppHandler {
	return &AppHandler{
		logger: generateLogger(name),
	}
}
