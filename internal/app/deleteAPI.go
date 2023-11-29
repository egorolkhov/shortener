package app

import (
	"encoding/json"
	"github.com/egorolkhov/shortener/internal/middleware"
	"github.com/egorolkhov/shortener/internal/storage"
	"io"
	"log"
	"net/http"
)

func (a *App) DeleteAPI(w http.ResponseWriter, r *http.Request) {
	responseData, err := io.ReadAll(r.Body)
	defer r.Body.Close() //закрывать все тела запроса
	if err != nil {
		log.Fatal(err)
	}
	var codes []string
	err = json.Unmarshal(responseData, &codes)

	cookie := w.Header().Get("Authorization")
	userID := middleware.GetUserID(cookie, "1234")

	err = storage.DeleteURL(r.Context(), a.DatabaseDSN, userID, codes)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusAccepted)
}
