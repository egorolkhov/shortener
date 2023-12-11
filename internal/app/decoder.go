package app

import (
	"github.com/egorolkhov/shortener/internal/storage"
	"github.com/gorilla/mux"
	"net/http"
)

func (a *App) DecodeURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var url string
	var err error

	if db, ok := a.Storage.(*storage.DB); ok {
		if storage.IsDeleted(r.Context(), db.DB, id) {
			w.WriteHeader(http.StatusGone)
			return
		}
	}

	url, err = a.Storage.Get(r.Context(), id)

	if err != nil {
		http.Error(w, "error when getting from storage", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err = w.Write([]byte("Location: " + url))
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}
