package storage

import (
	"errors"
	"sync"
)

type Data struct {
	sync.RWMutex
	Urls map[string]string
}

func New() *Data {
	return &Data{sync.RWMutex{}, make(map[string]string)}
}

func (u *Data) Add(short, fullurl string) {
	u.Lock()
	defer u.Unlock()
	u.Urls[short] = fullurl
}

func (u *Data) Get(short string) (string, error) {
	u.Lock()
	defer u.Unlock()
	_, ok := u.Urls[short]
	if !ok {
		return "", errors.New("no such id")
	}
	return u.Urls[short], nil
}
