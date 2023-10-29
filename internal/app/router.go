package app

import (
	config "github.com/egorolkhov/shortener/internal/config"
	"github.com/egorolkhov/shortener/internal/logger"
	"github.com/gorilla/mux"
)

func NewRouter(cfg *config.Cfg) *mux.Router {
	app := New(cfg)

	router := mux.NewRouter()

	router.HandleFunc("/", logger.PostLogger(app.ShortURL)).Methods("POST")
	router.HandleFunc("/{id}", logger.GetLogger(app.DecodeURL)).Methods("GET")
	router.HandleFunc("/api/shorten", logger.PostLogger(app.ShortAPI))
	return router
}
