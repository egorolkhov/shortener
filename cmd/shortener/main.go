package main

import (
	"context"
	"github.com/egorolkhov/shortener/internal/app"
	"github.com/egorolkhov/shortener/internal/config"
	"github.com/egorolkhov/shortener/internal/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	cfg := config.Config()

	r := app.NewRouter(cfg)

	srv := http.Server{
		Addr:    cfg.Address.String(),
		Handler: r,
	}

	err := logger.InitLogger()
	if err != nil {
		log.Println(err)
	}

	done := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt)
		<-sigs
		if err = srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(done)
	}()

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
