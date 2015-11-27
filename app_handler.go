package altserv

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

	"github.com/i2bskn/altserv/converter"
)

const errorTemplate = `
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

type errorInfo struct {
	Code    int
	Message string
}

func newErrorInfo(code int, message string) *errorInfo {
	return &errorInfo{
		Code:    code,
		Message: message,
	}
}

type ContentCache struct {
	Content     []byte
	ContentType string
	CachedAt    time.Time
}

func newContentCache(content []byte, contentType string) *ContentCache {
	return &ContentCache{
		Content:     content,
		ContentType: contentType,
		CachedAt:    time.Now(),
	}
}

type AppHandler struct {
	Config     *Config
	Converters *converter.AvailableConverters
	Caches     map[string]*ContentCache
}

func NewAppHandler(config *Config) *AppHandler {
	return &AppHandler{
		Config:     config,
		Converters: converter.NewAvailableConverters(),
		Caches:     make(map[string]*ContentCache),
	}
}

func (h *AppHandler) AccessLog(r *http.Request, status int) {
	logInfo := []string{
		r.Method,
		r.URL.Path,
		fmt.Sprint(status),
	}

	h.Config.Logger.Println(strings.Join(logInfo, " "))
}

func (h *AppHandler) RenderContent(w http.ResponseWriter, r *http.Request, content []byte, contentType string) {
	h.AccessLog(r, http.StatusOK)

	if len(contentType) == 0 {
		contentType = http.DetectContentType(content)
	}

	w.Header().Add("Content-Type", contentType)
	w.Write(content)
}

func (h *AppHandler) RenderError(w http.ResponseWriter, r *http.Request, i *errorInfo) {
	h.AccessLog(r, i.Code)
	t, _ := template.New("error").Parse(errorTemplate)
	w.WriteHeader(i.Code)
	t.Execute(w, i)
}

func (h *AppHandler) AssetPath(uri string) (string, os.FileInfo, error) {
	asset := filepath.Join(h.Config.DocumentRoot, uri)
	assetInfo, err := os.Stat(asset)

	if err == nil {
		if assetInfo.IsDir() {
			asset = filepath.Join(asset, h.Config.Index)
			assetInfo, err = os.Stat(asset)
			if err == nil {
				return asset, assetInfo, nil
			}
		} else {
			return asset, assetInfo, nil
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
			assetInfo, err = os.Stat(candidate)
			if err == nil {
				return candidate, assetInfo, nil
			}
		}
	}

	return "", nil, errors.New("File not found: " + asset)
}

func (h *AppHandler) Convert(src []byte, asset string) ([]byte, string) {
	ext := filepath.Ext(asset)
	content, toExt := h.Converters.Convert(src, ext)

	var contentType string
	switch toExt {
	case ".html":
		contentType = "text/html; charset=utf-8"
	case ".css":
		contentType = "text/css; charset=utf-8"
	}
	h.Caches[asset] = newContentCache(content, contentType)
	return content, contentType
}

func (h *AppHandler) ContentFromCache(asset string, info os.FileInfo) ([]byte, string) {
	contentCache, exist := h.Caches[asset]
	if exist {
		if contentCache.CachedAt.Sub(info.ModTime()) > 0 {
			return contentCache.Content, contentCache.ContentType
		}
		delete(h.Caches, asset)
	}
	return nil, ""
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	asset, info, err := h.AssetPath(r.URL.Path)
	if err != nil {
		info := newErrorInfo(http.StatusNotFound, err.Error())
		h.RenderError(w, r, info)
		return
	}

	cachedContent, cachedContentType := h.ContentFromCache(asset, info)
	if cachedContent != nil {
		h.RenderContent(w, r, cachedContent, cachedContentType)
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

	content, contentType := h.Convert(src, asset)
	h.RenderContent(w, r, content, contentType)
}
