package postgreSQLaddress

import (
	"errors"
	"strings"
)

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
	if len(res) != 5 {
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
		User:     "egoro",
		Password: "qazxsw",
		DBname:   "shortener",
		SSLmode:  "disable",
	}
}
