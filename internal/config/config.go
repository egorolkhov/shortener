package config

import (
	"flag"
	"github.com/egorolkhov/shortener/internal/config/netAddress"
	"os"
)

type Cfg struct {
	Address netaddress.NetAddress
	BaseURL string
}

func Config() *Cfg {
	//data := storage.New()
	address := netaddress.NewNetAddress()
	baseURL := ""

	flag.Var(address, "a", "http server adress")
	url := flag.String("b", baseURL, "base url address")

	flag.Parse()

	setURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		url = &setURL
	}
	_ = address.Set(os.Getenv("SERVER_ADDRESS"))

	return &Cfg{Address: *address, BaseURL: *url}
}
