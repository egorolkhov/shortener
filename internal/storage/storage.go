package storage

type Storage interface {
	Add(userID, code, url string) error
	Get(code string) (string, error)
	GetExist(url string) (string, error)
	GetUserURLS(userID string) []URL
}
