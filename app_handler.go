package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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

type ContentCache struct {
	Content  []byte
	CachedAt time.Time
}

func newContentCache(content []byte) *ContentCache {
	return &ContentCache{
		Content:  content,
		CachedAt: time.Now(),
	}
}

type AppHandler struct {
	Config     *Config
	Converters *AvailableConverters
	Caches     map[string]*ContentCache
}

func newAppHandler(config *Config) *AppHandler {
	return &AppHandler{
		Config:     config,
		Converters: newAvailableConverters(),
		Caches:     make(map[string]*ContentCache),
	}
}

func (h *AppHandler) AccessLog(r *http.Request, status int) {
	log_info := []string{
		r.Method,
		r.URL.Path,
		fmt.Sprint(status),
	}

	h.Config.Logger.Println(strings.Join(log_info, " "))
}

func (h *AppHandler) RenderContent(w http.ResponseWriter, r *http.Request, content []byte) {
	h.AccessLog(r, http.StatusOK)
	mime_type := http.DetectContentType(content)
	w.Header().Add("Content-Type", mime_type)
	w.Write(content)
}

func (h *AppHandler) RenderError(w http.ResponseWriter, r *http.Request, i *ErrorInfo) {
	h.AccessLog(r, i.Code)
	t, _ := template.New("error").Parse(ErrorTemplate)
	w.WriteHeader(i.Code)
	t.Execute(w, i)
}

func (h *AppHandler) AssetPath(uri string) (string, os.FileInfo, error) {
	asset := filepath.Join(h.Config.DocumentRoot, uri)
	asset_info, err := os.Stat(asset)

	if err == nil {
		if asset_info.IsDir() {
			asset = filepath.Join(asset, h.Config.Index)
			asset_info, err = os.Stat(asset)
			if err == nil {
				return asset, asset_info, nil
			}
		} else {
			return asset, asset_info, nil
		}
	}

	ext := filepath.Ext(asset)
	if len(ext) == 0 {
		ext = ".html"
	}

	if candidates, exist := h.Converters.ConvertMap[ext]; exist {
		dir, file := filepath.Split(asset)
		base := strings.TrimRight(file, ext)

		for _, c := range candidates {
			candidate := filepath.Join(dir, base+c)
			asset_info, err = os.Stat(candidate)
			if err == nil {
				return candidate, asset_info, nil
			}
		}
	}

	return "", nil, errors.New("File not found: " + asset)
}

func (h *AppHandler) Convert(src []byte, asset string) []byte {
	ext := filepath.Ext(asset)
	converted_src := h.Converters.Convert(src, ext)
	h.Caches[asset] = newContentCache(converted_src)
	return converted_src
}

func (h *AppHandler) ContentFromCache(asset string, info os.FileInfo) []byte {
	content_cache, exist := h.Caches[asset]
	if exist {
		if content_cache.CachedAt.Sub(info.ModTime()) > 0 {
			return content_cache.Content
		}
		delete(h.Caches, asset)
	}
	return nil
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	asset, info, err := h.AssetPath(r.URL.Path)
	if err != nil {
		info := newErrorInfo(http.StatusNotFound, err.Error())
		h.RenderError(w, r, info)
		return
	}

	content := h.ContentFromCache(asset, info)
	if content != nil {
		h.RenderContent(w, r, content)
		return
	}

	f, err := os.OpenFile(asset, os.O_RDONLY, 0)
	if err != nil {
		info := newErrorInfo(http.StatusInternalServerError, err.Error())
		h.RenderError(w, r, info)
		return
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		info := newErrorInfo(http.StatusInternalServerError, err.Error())
		h.RenderError(w, r, info)
		return
	}

	h.RenderContent(w, r, h.Convert(src, asset))
}
