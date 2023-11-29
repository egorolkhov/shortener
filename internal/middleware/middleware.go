package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func Middleware(h http.HandlerFunc) http.HandlerFunc {
	foo := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), `gzip`) {
			h.ServeHTTP(w, r)
			return
		}
		gz := gzip.NewWriter(w)
		defer gz.Close()
		cw := &CompressWrite{w, gz}
		cw.Header().Set("Content-Encoding", "gzip")

		if !strings.Contains(r.Header.Get("Content-Encoding"), `gzip`) {
			h.ServeHTTP(cw, r)
			return
		}

		gzR, err := gzip.NewReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cr := &CompressRead{r.Body, gzR}

		r.Body = cr

		h.ServeHTTP(cw, r)
	}
	return foo
}
