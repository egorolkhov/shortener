package app

import (
	"errors"
	"github.com/egorolkhov/shortener/internal/app/encoder"
	"github.com/egorolkhov/shortener/internal/middleware"
	"github.com/egorolkhov/shortener/internal/storage"
	"io"
	"log"
	"net/http"
)

func (a *App) ShortURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var temp int
	responseData, err := io.ReadAll(r.Body)
	defer r.Body.Close() //закрывать все тела запроса
	if err != nil {
		log.Fatal(err)
	}

	url := string(responseData)

	code := encoder.Code()

	w.Header().Set("Content-Type", "text/plain")

	cookie := w.Header().Get("Authorization")
	userID := middleware.GetUserID(cookie)
	if userID == "error" {
		w.WriteHeader(http.StatusUnauthorized)
	}
	err = a.Storage.Add(r.Context(), userID, code, url)
	if errors.Is(err, storage.ErrURLAlreadyExist) {
		temp = 1
		w.WriteHeader(http.StatusConflict)
		code, err = a.Storage.GetExist(r.Context(), url)
		if err != nil {
			log.Println(err)
		}
	}

	//fmt.Println(a.BaseURL)
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
