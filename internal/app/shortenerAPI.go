package app

import (
	"encoding/json"
	"fmt"
	"github.com/egorolkhov/shortener/internal/app/encoder"
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
	var url RequestData

	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer r.Body.Close() //закрывать все тела запроса

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	code := encoder.Code()

	a.Storage.Add(code, url.URL)

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
	w.Write(result)

	log.Println(a.Storage.Urls)
}
