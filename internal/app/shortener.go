package app

import (
	"errors"
	"fmt"
	"github.com/egorolkhov/shortener/internal/app/encoder"
	"github.com/egorolkhov/shortener/internal/storage"
	"io"
	"log"
	"net/http"
)

func (a *App) ShortURL(w http.ResponseWriter, r *http.Request) {
	var temp int
	responseData, err := io.ReadAll(r.Body)
	defer r.Body.Close() //закрывать все тела запроса
	if err != nil {
		log.Fatal(err)
	}

	url := string(responseData)

	code := encoder.Code()

	if a.flag == 1 {
		err = storage.AddDB(r.Context(), a.DatabaseDSN, code, url)
		if errors.Is(err, storage.ErrURLAlreadyExist) {
			temp = 1
			w.WriteHeader(http.StatusConflict)
			code, err = storage.GetDBExist(r.Context(), a.DatabaseDSN, url)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		err = a.Storage.Add(code, url)
		if errors.Is(err, storage.ErrURLAlreadyExist) {
			temp = 1
			w.WriteHeader(http.StatusConflict)
			code, err = a.Storage.GetExist(url)
			if err != nil {
				log.Println(err)
			}
		}
	}
	if err != nil {
		log.Println(err)
	}

	fmt.Println(a.BaseURL)
	var resp string
	if a.BaseURL != "" {
		resp = a.BaseURL + "/" + code
	} else {
		resp = "http://" + r.Host + "/" + a.BaseURL + code
	}

	err = storage.FileWrite(code, url, a.Filepath)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "text/plain")
	if temp != 1 {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(resp))

	log.Println(a.Storage.Urls)
}
