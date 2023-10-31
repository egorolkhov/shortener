package config

import (
	"flag"
	"github.com/egorolkhov/shortener/internal/config/netAddress"
	"os"
)

type Cfg struct {
	Address  netaddress.NetAddress
	BaseURL  string
	Filepath string
}

func Config() *Cfg {
	address := netaddress.NewNetAddress()
	baseURL := ""

	flag.Var(address, "a", "http server adress")
	url := flag.String("b", baseURL, "base url address")
	filepath := flag.String("f", "./tmp/short-url-db.json", "db filepath")

	flag.Parse()

	setURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		url = &setURL
	}
	setFilepath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		filepath = &setFilepath
	}
	_ = address.Set(os.Getenv("SERVER_ADDRESS"))

	return &Cfg{Address: *address, BaseURL: *url, Filepath: *filepath}
}
