package app

import (
	"fmt"
	"github.com/egorolkhov/shortener/internal/app/encoder"
	"io"
	"log"
	"net/http"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (a *App) ShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(400)
		return
	}

	responseData, err := io.ReadAll(r.Body)
	defer r.Body.Close() //закрывать все тела запроса
	if err != nil {
		log.Fatal(err)
	}

	url := string(responseData)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(201)

	code := encoder.Code()

	a.Storage.Add(code, url)

	fmt.Println(a.BaseURL)
	var resp string
	if a.BaseURL != "" {
		resp = a.BaseURL + "/" + code
	} else {
		resp = "http://" + r.Host + "/" + a.BaseURL + code
	}
	//resp := "http://" + r.Host + "/" + a.BaseURL + code
	//resp := a.BaseURL + "/" + code

	w.Write([]byte(resp))
	//fmt.Println(a.Storage)
}
