package app

import (
	config "github.com/egorolkhov/shortener/internal/config"
	"github.com/gorilla/mux"
)

func NewRouter(cfg *config.Cfg) *mux.Router {
	var app Handler

	app = New(cfg)

	router := mux.NewRouter()
	router.HandleFunc("/", app.ShortURL).Methods("POST")
	router.HandleFunc("/{id}", app.DecodeURL).Methods("GET")
	return router
}
