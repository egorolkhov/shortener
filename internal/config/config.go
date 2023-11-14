package config

import (
	"flag"
	"github.com/egorolkhov/shortener/internal/config/netAddress"
	"os"
)

type Cfg struct {
	Address     netAddress.NetAddress
	BaseURL     string
	Filepath    string
	DatabaseDSN string
}

func Config() *Cfg {
	address := netAddress.NewNetAddress()
	baseURL := ""

	flag.Var(address, "a", "http server adress")
	url := flag.String("b", baseURL, "base url address")
	filepath := flag.String("f", "", "db filepath")
	//databaseDSN := flag.String("d", "host=localhost port=5432 user=shortener password=Qazxsw2200 dbname=shortener sslmode=disable", "db address")
	databaseDSN := flag.String("d", "", "db address")
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

	return &Cfg{Address: *address, BaseURL: *url, Filepath: *filepath, DatabaseDSN: *databaseDSN}
}
