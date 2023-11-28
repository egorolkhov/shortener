package storage

import (
	"errors"
	"fmt"
	"sync"
)

type Data struct {
	mu    sync.RWMutex
	Urls  map[string]string
	Codes map[string]string
}

func NewLocalData() *Data {
	return &Data{sync.RWMutex{}, make(map[string]string), make(map[string]string)}
}

func (u *Data) Add(code, url string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.Urls[url]; ok {
		fmt.Println(u.Urls)
		return ErrURLAlreadyExist
	}
	u.Urls[url] = code
	u.Codes[code] = url
	return nil
}

func (u *Data) Get(code string) (string, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	if _, ok := u.Codes[code]; !ok {
		return "", errors.New("no such id")
	}
	return u.Codes[code], nil
}

func (u *Data) GetExist(url string) (string, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	if v, ok := u.Urls[url]; ok {
		return v, nil
	}
	return "", errors.New("error when getting short url")
}
