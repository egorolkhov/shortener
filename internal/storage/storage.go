package storage

import (
	"errors"
	"sync"
)

type Data struct {
	mu   sync.RWMutex
	Urls map[string]string
}

func New() *Data {
	return &Data{sync.RWMutex{}, make(map[string]string)}
}

func (u *Data) Add(short, fullURL string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, v := range u.Urls {
		if v == fullURL {
			return ErrURLAlreadyExist
		}
	}
	u.Urls[short] = fullURL
	return nil
}

func (u *Data) Get(short string) (string, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	_, ok := u.Urls[short]
	if !ok {
		return "", errors.New("no such id")
	}
	return u.Urls[short], nil
}

func (u *Data) GetExist(fullURL string) (string, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	for k, v := range u.Urls {
		if v == fullURL {
			return k, nil
		}
	}
	return "", errors.New("error when getting short url")
}
