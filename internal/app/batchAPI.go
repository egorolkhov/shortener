package app

import (
	"encoding/json"
	"github.com/egorolkhov/shortener/internal/app/encoder"
	"github.com/egorolkhov/shortener/internal/storage"
	"log"
	"net/http"
)

func (a *App) BatchAPI(w http.ResponseWriter, r *http.Request) {
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

	err = storage.AddBatch(r.Context(), a.DatabaseDSN, codes, jsons)
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
			a.Storage.Add(codes[i], json.OriginalURL)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)

	log.Println(a.Storage.Urls)
}
