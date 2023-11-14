package config

import (
	"errors"
	"flag"
	"github.com/egorolkhov/shortener/internal/config/netAddress"
	"os"
	"strings"
)

type Cfg struct {
	Address     netAddress.NetAddress
	BaseURL     string
	Filepath    string
	DatabaseDSN PGXaddress
}

func Config() *Cfg {
	address := netAddress.NewNetAddress()
	databaseDSN := NewPGXaddress()
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

type PGXaddress struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
	SSLmode  string
}

func (n *PGXaddress) String() string {
	sb := strings.Builder{}

	sb.WriteString(n.Host)
	sb.WriteString(":")
	sb.WriteString(n.Port)
	sb.WriteString(":")
	sb.WriteString(n.User)
	sb.WriteString(":")
	sb.WriteString(n.Password)
	sb.WriteString(":")
	sb.WriteString(n.DBname)
	sb.WriteString(":")
	sb.WriteString(n.SSLmode)

	return sb.String()
}

func (n *PGXaddress) Set(flagValue string) error {
	res := strings.Split(flagValue, ":")
	if len(res) != 6 {
		return errors.New("wrong address")
	}
	n.Host = res[0]
	n.Port = res[1]
	n.User = res[2]
	n.Password = res[3]
	n.DBname = res[4]
	n.SSLmode = res[5]
	return nil
}

func NewPGXaddress() *PGXaddress {
	return &PGXaddress{
		Host:     "localhost",
		Port:     "5432",
		User:     "shortener",
		Password: "Qazxsw2200",
		DBname:   "shortener",
		SSLmode:  "disable",
	}
}
