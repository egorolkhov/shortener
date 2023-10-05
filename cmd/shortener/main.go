package main

import (
	"github.com/egorolkhov/shortener/internal/app"
	"github.com/egorolkhov/shortener/internal/config"
	"net/http"
)

//iter3

func main() {
	cfg := config.Config()

	r := app.NewRouter(cfg)

	err := http.ListenAndServe(cfg.Address.String(), r)
	if err != nil {
		panic(err)
	}
}
