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

	router.HandleFunc("/api/user/urls", middleware.Cookie(app.DeleteAPI))
	router.HandleFunc("/ping", middleware.Cookie(middleware.Middleware(logger.GetLogger(app.PSQLconnection)))).Methods("GET")
	router.HandleFunc("/", middleware.Cookie(middleware.Middleware(logger.PostLogger(app.ShortURL)))).Methods("POST")
	router.HandleFunc("/{id}", middleware.Cookie(middleware.Middleware(logger.GetLogger(app.DecodeURL)))).Methods("GET")
	router.HandleFunc("/api/shorten", middleware.Cookie(middleware.Middleware(logger.PostLogger(app.ShortAPI)))).Methods("POST")
	router.HandleFunc("/api/shorten/batch", middleware.Cookie(middleware.Middleware(logger.PostLogger(app.BatchAPI)))).Methods("POST")
	router.HandleFunc("/api/user/urls", middleware.Cookie(middleware.Middleware(logger.PostLogger(app.UserAPI)))).Methods("GET")
	return router
}
