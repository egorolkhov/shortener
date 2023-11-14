package config

import (
	"flag"
	"github.com/egorolkhov/shortener/internal/config/netAddress"
	"github.com/egorolkhov/shortener/internal/config/postgresqladdress"
	"os"
)

type Cfg struct {
	Address     netAddress.NetAddress
	BaseURL     string
	Filepath    string
	DatabaseDSN postgresqladdress.PGXaddress
}

func Config() *Cfg {
	address := netAddress.NewNetAddress()
	databaseDSN := postgresqladdress.NewPGXaddress()
	baseURL := ""

	flag.Var(address, "a", "http server adress")
	url := flag.String("b", baseURL, "base url address")
	filepath := flag.String("f", "./tmp/short-url-db.json", "db filepath")
	flag.Var(databaseDSN, "d", "db address")
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
	_ = databaseDSN.Set(os.Getenv("DATABASE_DSN"))

	return &Cfg{Address: *address, BaseURL: *url, Filepath: *filepath, DatabaseDSN: *databaseDSN}
}
