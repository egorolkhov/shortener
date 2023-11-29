package middleware

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func Cookie(h http.HandlerFunc) http.HandlerFunc {
	foo := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("default")
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
					Expires: time.Now().Add(24 * time.Hour)}
				http.SetCookie(w, cookie)
			default:
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
		w.Header().Set("Authorization", cookie.Value)

		h.ServeHTTP(w, r)
	}
	return foo
}
