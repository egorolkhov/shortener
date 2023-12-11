package middleware

import (
	"bufio"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

const secretKeyPath = "tmp/key.txt"

var secretKey, _ = getKey(secretKeyPath)

func Cookie(h http.HandlerFunc) http.HandlerFunc {
	foo := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("default")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				value, err := BuidToken(secretKey)
				if err != nil {
					log.Println(err)
				}
				cookie = &http.Cookie{
					Name:    "default",
					Value:   value,
					Expires: time.Now().Add(TokenExp)}
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

func getKey(filepath string) (string, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var line string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}
	return line, nil
}
