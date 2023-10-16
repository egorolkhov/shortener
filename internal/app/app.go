package app

import (
	"github.com/egorolkhov/shortener/internal/config"
	storage "github.com/egorolkhov/shortener/internal/storage"
	"net/http"
)

type App struct {
	Storage *storage.Data
	BaseURL string
}

type Handler interface {
	DecodeURL(w http.ResponseWriter, r *http.Request)
	ShortURL(w http.ResponseWriter, r *http.Request)
}

func New(cfg *config.Cfg) *App {
	return &App{storage.New(), cfg.BaseURL}
}
