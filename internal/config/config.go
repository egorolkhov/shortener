package config

import (
	"flag"
	"fmt"
	"github.com/egorolkhov/shortener/internal/storage"
	"os"
)

type Cfg struct {
	Address NetAddress
	BaseURL string
	Data    *storage.Data
}

func Config() *Cfg {
	data := storage.New()
	address := NewNetAddress()
	baseURL := NewBaseURL()

	flag.Var(address, "a", "http server adress")
	url := flag.String("b", baseURL, "base url address")

	flag.Parse()

	setURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		url = &setURL
	}
	err := address.Set(os.Getenv("SERVER_ADDRESS"))
	if err != nil {
		fmt.Println("wrong address format")
	}

	return &Cfg{Address: *address, BaseURL: *url, Data: data}
}
