package request

import (
	"errors"
	"mime"
	"net/http"
	"path/filepath"
)

// FilenameFromResponse analysis filename from http response
func FilenameFromResponse(resp *http.Response) (filename string, err error) {
	_, info, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	filename = info["filename"]
	if err == nil && filename != "" {
		return
	}

	filename = filepath.Base(resp.Request.URL.Path)
	if filename != "." {
		err = nil
		return
	}

	err = errors.New("filename not found in response")
	return
}
