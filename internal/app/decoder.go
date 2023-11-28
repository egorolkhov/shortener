package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (a *App) DecodeURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var url string
	var err error

	url, err = a.Storage.Get(id)
	//if a.flag == 1 {
	//	url, err = storage.GetDB(r.Context(), a.DatabaseDSN, id)
	//} else {
	//	url, err = a.Storage.Get(id)
	//}
	//url, err = storage.GetDB(r.Context(), a.DatabaseDSN, id)
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
