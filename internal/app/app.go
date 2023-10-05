package app

import (
	"github.com/egorolkhov/shortener/internal/config"
	storage "github.com/egorolkhov/shortener/internal/storage"
)

type App struct {
	Storage *storage.Data
	BaseURL string
}

func New(cfg *config.Cfg) *App {
	return &App{cfg.Data, cfg.BaseURL}
}
