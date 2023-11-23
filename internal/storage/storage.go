package storage

type Storage interface {
	Add(code, url string) error
	Get(code string) (string, error)
	GetExist(url string) (string, error)
}
