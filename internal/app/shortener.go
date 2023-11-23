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

	w.Header().Set("Content-Type", "text/plain")

	err = a.Storage.Add(code, url)
	if errors.Is(err, storage.ErrURLAlreadyExist) {
		temp = 1
		w.WriteHeader(http.StatusConflict)
		code, err = a.Storage.GetExist(url)
		if err != nil {
			log.Println(err)
		}
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

	if temp != 1 {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(resp))

	if storage, ok := a.Storage.(*storage.Data); ok {
		log.Println(storage.Urls)
		log.Println(storage.Codes)
	}
	//  log.Println(a.Storage.Urls)
}
