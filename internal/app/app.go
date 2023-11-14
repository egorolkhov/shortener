package app

import (
	"github.com/egorolkhov/shortener/internal/config"
	storage "github.com/egorolkhov/shortener/internal/storage"
	"log"
	"net/http"
)

type App struct {
	Storage  *storage.Data
	BaseURL  string
	Filepath string
	//DatabaseDSN config.PGXaddress
	DatabaseDSN string
	flag        int
}

type Handler interface {
	DecodeURL(w http.ResponseWriter, r *http.Request)
	ShortURL(w http.ResponseWriter, r *http.Request)
}

func New(cfg *config.Cfg) *App {
	var flag int

	if cfg.DatabaseDSN != "" {
		err := storage.CreateTable(cfg.DatabaseDSN)
		if err != nil {
			log.Println(err)
		}
		flag = 1
	}

	Storage := storage.New()
	if cfg.Filepath != "" && flag != 1 {
		err := storage.GetStorage(Storage, cfg.Filepath)
		if err != nil {
			log.Println(err)
		}
	}

	return &App{
		Storage,
		cfg.BaseURL,
		cfg.Filepath,
		cfg.DatabaseDSN,
		flag}
}
