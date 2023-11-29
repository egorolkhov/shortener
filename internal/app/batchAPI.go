package app

import (
	"encoding/json"
	"fmt"
	"github.com/egorolkhov/shortener/internal/app/encoder"
	"github.com/egorolkhov/shortener/internal/middleware"
	"github.com/egorolkhov/shortener/internal/storage"
	"log"
	"net/http"
)

func (a *App) BatchAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ADD_BATCH")
	var jsons []storage.RequestBatch

	err := json.NewDecoder(r.Body).Decode(&jsons)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	var codes []string
	for i := 0; i < len(jsons); i++ {
		codes = append(codes, encoder.Code())
	}

	cookie := w.Header().Get("Authorization")
	userID := middleware.GetUserID(cookie, "1234")
	if userID == "error" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = storage.AddBatch(r.Context(), a.DatabaseDSN, userID, codes, jsons)
	if err != nil {
		log.Println(err)
	}

	var resp []storage.ResponseBatch

	for i := 0; i < len(jsons); i++ {
		var shortURL string
		if err != nil {
			log.Println(err)
		}
		if a.BaseURL != "" {
			shortURL = a.BaseURL + "/" + codes[i]
		} else {
			shortURL = "http://" + r.Host + "/" + a.BaseURL + codes[i]
		}
		resp = append(resp, storage.ResponseBatch{
			CorrelationID: jsons[i].CorrelationID,
			ShortURL:      shortURL})
	}

	result, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, json := range jsons {
		storage.FileWrite(codes[i], json.OriginalURL, a.Filepath)
		if a.flag != 1 {
			a.Storage.Add(userID, codes[i], json.OriginalURL)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)

	if storage, ok := a.Storage.(*storage.Data); ok {
		fmt.Println(storage.Urls)
		fmt.Println(storage.Codes)
		fmt.Println(storage.Users)
	}

	//log.Println(a.Storage.Urls)
}
