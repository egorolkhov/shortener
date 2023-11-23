package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/egorolkhov/shortener/internal/app/encoder"
	"github.com/egorolkhov/shortener/internal/storage"
	"log"
	"net/http"
)

type ResponseData struct {
	Result string `json:"result"`
}

type RequestData struct {
	URL string `json:"url"`
}

func (a *App) ShortAPI(w http.ResponseWriter, r *http.Request) {
	var temp int
	var url RequestData

	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer r.Body.Close() //закрывать все тела запроса

	code := encoder.Code()

	w.Header().Set("Content-Type", "application/json")

	err = a.Storage.Add(code, url.URL)
	if errors.Is(err, storage.ErrURLAlreadyExist) {
		temp = 1
		w.WriteHeader(http.StatusConflict)
		code, err = a.Storage.GetExist(url.URL)
		if err != nil {
			log.Println(err)
		}
	}

	if err != nil {
		log.Println(err)
	}

	var resp ResponseData
	if a.BaseURL != "" {
		resp.Result = a.BaseURL + "/" + code
	} else {
		resp.Result = "http://" + r.Host + "/" + a.BaseURL + code
	}

	result, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storage.FileWrite(code, url.URL, a.Filepath)

	if temp != 1 {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write(result)

	if storage, ok := a.Storage.(*storage.Data); ok {
		log.Println(storage.Urls)
		log.Println(storage.Codes)
	}
	//log.Println(a.Storage.Urls)
}
