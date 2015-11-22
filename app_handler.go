package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const ErrorTemplate = `
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Error</title>
</head>
<body>
<h1>Status {{.Code}}</h1>
<p>{{.Message}}</p>
</body>
</html>
`

type AppHandler struct {
	Config *Config
}

func newAppHandler(config *Config) *AppHandler {
	return &AppHandler{
		Config: config,
	}
}

type ErrorInfo struct {
	Code    int
	Message string
}

func newErrorInfo(code int, message string) *ErrorInfo {
	return &ErrorInfo{
		Code:    code,
		Message: message,
	}
}

func (h *AppHandler) RequestLog(r *http.Request) {
	log_info := []string{
		r.Method,
		r.URL.Path,
	}

	h.Config.Logger.Println(strings.Join(log_info, " "))
}

func (h *AppHandler) RenderError(w http.ResponseWriter, i *ErrorInfo) {
	t, _ := template.New("error").Parse(ErrorTemplate)
	w.WriteHeader(i.Code)
	t.Execute(w, i)
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.RequestLog(r)

	asset, err := h.Config.AssetPath(r.URL.Path)
	if err != nil {
		info := newErrorInfo(http.StatusNotFound, err.Error())
		h.RenderError(w, info)
		return
	}

	f, err := os.OpenFile(asset, os.O_RDONLY, 0)
	if err != nil {
		info := newErrorInfo(http.StatusInternalServerError, err.Error())
		h.RenderError(w, info)
		return
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		info := newErrorInfo(http.StatusInternalServerError, err.Error())
		h.RenderError(w, info)
		return
	}

	mime_type := http.DetectContentType(bytes)
	w.Header().Add("Content-Type", mime_type)
	w.Write(bytes)
}
