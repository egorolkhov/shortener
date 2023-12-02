package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Data struct {
	mu    sync.RWMutex
	Urls  map[string]string
	Codes map[string]string
	Users map[string][]URL
}

func NewLocalData() *Data {
	return &Data{sync.RWMutex{}, make(map[string]string), make(map[string]string), make(map[string][]URL)}
}

func (u *Data) Add(ctx context.Context, userID, code, url string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.Urls[url]; ok {
		fmt.Println(u.Urls)
		return ErrURLAlreadyExist
	}
	u.Urls[url] = code
	u.Codes[code] = url
	if _, ok := u.Codes[url]; !ok {
		u.Users[userID] = append(u.Users[userID], URL{FullURL: url, ShortURL: code})
	}
	//u.Users[userID] = append(u.Users[userID], URL{FullURL: url, ShortURL: code})
	return nil
}

func (u *Data) Get(ctx context.Context, code string) (string, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	if _, ok := u.Codes[code]; !ok {
		return "", errors.New("no such id")
	}
	return u.Codes[code], nil
}

func (u *Data) GetExist(ctx context.Context, url string) (string, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	if v, ok := u.Urls[url]; ok {
		return v, nil
	}
	return "", errors.New("error when getting short url")
}

func (u *Data) GetUserURLS(ctx context.Context, userID string) []URL {
	return u.Users[userID]
}
