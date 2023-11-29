package config

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const secretKeyPath = "tmp/key.txt"

type Cfg struct {
	Address     NetAddress
	BaseURL     string
	Filepath    string
	DatabaseDSN string
	SecretKey   string
}

func Config() *Cfg {
	address := NewNetAddress()
	baseURL := ""
	secretKey, err := getKey(secretKeyPath)
	if err != nil {
		fmt.Println("can't get secretKey")
	}

	flag.Var(address, "a", "http server adress")
	url := flag.String("b", baseURL, "base url address")
	filepath := flag.String("f", "", "db filepath")
	databaseDSN := flag.String("d", "host=localhost port=5432 user=shortener password=Qazxsw2200 dbname=shortener sslmode=disable", "db address")
	//databaseDSN := flag.String("d", "", "db address")
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

	//fmt.Println(*address, *url, *filepath, *databaseDSN)

	return &Cfg{Address: *address, BaseURL: *url, Filepath: *filepath, DatabaseDSN: *databaseDSN, SecretKey: secretKey}
}

func getKey(filepath string) (string, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var line string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}
	return line, nil
}

type NetAddress struct {
	Host string
	Port string
}

func (n *NetAddress) String() string {
	sb := strings.Builder{}

	sb.WriteString(n.Host)
	sb.WriteString(":")
	sb.WriteString(n.Port)

	return sb.String()
}

func (n *NetAddress) Set(flagValue string) error {
	res := strings.Split(flagValue, ":")
	if len(res) != 2 {
		return errors.New("wrong address")
	}
	n.Host = res[0]
	n.Port = res[1]
	return nil
}

func NewNetAddress() *NetAddress {
	return &NetAddress{Host: "localhost", Port: "8080"}
}
