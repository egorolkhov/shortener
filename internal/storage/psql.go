package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log"
)

var ErrURLAlreadyExist = errors.New("this URL already in database")

type URL struct {
	fullURL  string
	shortURL string
}

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

func AddDB(ctx context.Context, DatabaseDSN string, code, url string) error {
	db, err := sql.Open("pgx", DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	uuid := generateCorrelationID()

	result, err := db.ExecContext(ctx, "INSERT INTO original_urls (correlation_id, full_url)"+
		" VALUES ($1,$2) ON CONFLICT (full_url) DO NOTHING;", uuid, url)

	if i, _ := result.RowsAffected(); i == 0 {
		return ErrURLAlreadyExist
	}
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, "INSERT INTO short_url (correlation_id, short_url) "+
		" VALUES ($1,$2)", uuid, code)
	if err != nil {
		return err
	}
	return nil
}

func GetDB(ctx context.Context, DatabaseDSN string, short string) (string, error) {
	db, err := sql.Open("pgx", DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	var uuid string

	row := db.QueryRowContext(ctx, "SELECT correlation_id FROM short_url WHERE short_url = $1", short)
	err = row.Scan(&uuid)
	if err != nil {
		log.Println(err)
	}

	var url URL

	row = db.QueryRowContext(ctx, "SELECT full_url FROM original_urls WHERE correlation_id = $1", uuid)
	err = row.Scan(&url.fullURL)
	if err != nil {
		log.Println(err)
	}

	return url.fullURL, nil
}

func GetDBExist(ctx context.Context, DatabaseDSN string, url string) (string, error) {
	db, err := sql.Open("pgx", DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	var id string

	row := db.QueryRowContext(ctx, "SELECT correlation_id FROM original_urls WHERE full_url = $1", url)
	err = row.Scan(&id)
	if err != nil {
		log.Println(err)
	}

	var short string

	row = db.QueryRowContext(ctx, "SELECT short_url FROM short_url WHERE correlation_id = $1", id)
	err = row.Scan(&short)
	if err != nil {
		log.Println(err)
	}

	return short, nil
}

func AddBatch(ctx context.Context, DatabaseDSN string, codes []string, jsons []RequestBatch) error {
	db, err := sql.Open("pgx", DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmtShort, err := tx.PrepareContext(ctx, "INSERT INTO short_url (correlation_id, short_url)"+
		" VALUES ($1,$2)")
	if err != nil {
		return err
	}
	defer stmtShort.Close()

	stmtOriginal, err := tx.PrepareContext(ctx, "INSERT INTO original_urls (correlation_id, full_url) "+
		"VALUES($1,$2)")
	if err != nil {
		return err
	}
	defer stmtOriginal.Close()

	for i, line := range jsons {
		_, err = stmtShort.ExecContext(ctx, line.CorrelationID, codes[i])
		if err != nil {
			return err
		}
		_, err = stmtOriginal.ExecContext(ctx, line.CorrelationID, line.OriginalURL)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
