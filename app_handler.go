package main

import (
	"io"
	"net/http"
	"strings"
)

type AppHandler struct {
	Config *Config
}

func (h *AppHandler) RequestLog(r *http.Request) {
	log_info := []string{
		r.Method,
		r.URL.Path,
	}

	h.Config.Logger.Println(strings.Join(log_info, " "))
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.RequestLog(r)
	asset, _ := h.Config.AssetPath(r.URL.Path)
	io.WriteString(w, asset)
}

func newAppHandler(config *Config) *AppHandler {
	return &AppHandler{
		Config: config,
	}
}
