package middleware

import (
	"compress/gzip"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func Middleware(h http.HandlerFunc) http.HandlerFunc {
	foo := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("default")
		fmt.Println(err)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				value, err := BuidToken("1234")
				if err != nil {
					log.Println(err)
				}
				cookie = &http.Cookie{
					Name:    "default",
					Value:   value,
					Expires: time.Now().Add(365 * 24 * time.Hour)}
				http.SetCookie(w, cookie)
			default:
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				//return
				//w.WriteHeader(http.StatusInternalServerError)
			}
		}
		w.Header().Set("Authorization", cookie.Value)

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
