package storage

import "context"

type Storage interface {
	Add(ctx context.Context, userID, code, url string) error
	Get(ctx context.Context, code string) (string, error)
	GetExist(ctx context.Context, url string) (string, error)
	GetUserURLS(ctx context.Context, userID string) []URL
}
