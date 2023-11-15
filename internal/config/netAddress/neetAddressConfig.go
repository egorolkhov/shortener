package netAddress

import (
	"errors"
	"strings"
)

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
