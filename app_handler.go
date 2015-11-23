package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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
	Config     *Config
	Converters *AvailableConverters
}

func newAppHandler(config *Config) *AppHandler {
	return &AppHandler{
		Config:     config,
		Converters: newAvailableConverters(),
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

func (h *AppHandler) AssetPath(uri string) (asset string, err error) {
	var asset_info os.FileInfo
	asset = filepath.Join(h.Config.DocumentRoot, uri)
	asset_info, err = os.Stat(asset)

	if err == nil {
		if asset_info.IsDir() {
			asset = filepath.Join(asset, h.Config.Index)
			asset_info, err = os.Stat(asset)
			if err == nil {
				return asset, nil
			}
		} else {
			return asset, nil
		}
	}

	ext := filepath.Ext(asset)
	if len(ext) == 0 {
		ext = ".html"
	}

	candidates, exist := h.Converters.ConvertMap[ext]
	if exist {
		dir, file := filepath.Split(asset)
		base := strings.TrimRight(file, ext)

		for _, c := range candidates {
			candidate := filepath.Join(dir, base+c)
			asset_info, err = os.Stat(candidate)
			if err == nil {
				return candidate, nil
			}
		}
	}

	return "", errors.New("File not found: " + asset)
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.RequestLog(r)

	asset, err := h.AssetPath(r.URL.Path)
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

	asset_bytes, err := ioutil.ReadAll(f)
	if err != nil {
		info := newErrorInfo(http.StatusInternalServerError, err.Error())
		h.RenderError(w, info)
		return
	}

	converted_bytes := h.Converters.Convert(asset_bytes, filepath.Ext(asset))
	mime_type := http.DetectContentType(converted_bytes)
	w.Header().Add("Content-Type", mime_type)
	w.Write(converted_bytes)
}
