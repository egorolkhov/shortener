package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (a *App) DecodeURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "path doesn`t exist", http.StatusNotFound)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	url, err := a.Storage.Get(id)
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
