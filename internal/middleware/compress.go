package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
)

var CompresTypes = []string{
	"application/json",
	"text/html",
}

type CompressWrite struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func (c *CompressWrite) Write(b []byte) (int, error) {
	res, err := c.zw.Write(b)
	return res, err
}

type CompressRead struct {
	io.ReadCloser
	zr *gzip.Reader
}

func (c *CompressRead) Read(p []byte) (int, error) {
	res, err := c.zr.Read(p)
	return res, err
}
