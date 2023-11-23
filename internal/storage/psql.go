package storage

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log"
)

var ErrURLAlreadyExist = errors.New("this URL already in database")

type RequestBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func generateCorrelationID() string {
	id := uuid.New()
	return id.String()
}

func CreateTable(DatabaseDSN string) error {
	db, err := sql.Open("pgx", DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	TableShortURL := `
        CREATE TABLE IF NOT EXISTS short_url (
            correlation_id TEXT,
            short_url TEXT UNIQUE
        )
    `

	TableOriginalURL := `
        CREATE TABLE IF NOT EXISTS original_urls (
            correlation_id TEXT,
            full_url TEXT UNIQUE
        )
    `

	_, err = db.Exec(TableShortURL)
	if err != nil {
		return err
	}

	_, err = db.Exec(TableOriginalURL)
	if err != nil {
		return err
	}
	return nil
}
