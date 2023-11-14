package app

import (
	config "github.com/egorolkhov/shortener/internal/config"
	"github.com/egorolkhov/shortener/internal/logger"
	"github.com/egorolkhov/shortener/internal/middleware"
	"github.com/gorilla/mux"
)

func NewRouter(cfg *config.Cfg) *mux.Router {
	app := New(cfg)

	router := mux.NewRouter()

	router.HandleFunc("/ping", middleware.Middleware(logger.GetLogger(app.PSQLconnection))).Methods("GET")
	router.HandleFunc("/", middleware.Middleware(logger.PostLogger(app.ShortURL))).Methods("POST")
	router.HandleFunc("/{id}", middleware.Middleware(logger.GetLogger(app.DecodeURL))).Methods("GET")
	router.HandleFunc("/api/shorten", middleware.Middleware(logger.PostLogger(app.ShortAPI))).Methods("POST")
	return router
}
