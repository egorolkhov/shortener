package storage

import (
	"context"
	"database/sql"
	"log"
)

type URL struct {
	fullURL  string
	shortURL string
}

type DB struct {
	ctx         context.Context
	DatabaseDSN string
}

func NewDB(DatabaseDSN string) *DB {
	return &DB{context.Background(), DatabaseDSN}
}

func (d *DB) Add(code, url string) error {
	db, err := sql.Open("pgx", d.DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	uuid := generateCorrelationID()

	result, err := db.ExecContext(d.ctx, "INSERT INTO original_urls (correlation_id, full_url)"+
		" VALUES ($1,$2) ON CONFLICT (full_url) DO NOTHING;", uuid, url)

	if i, _ := result.RowsAffected(); i == 0 {
		return ErrURLAlreadyExist
	}
	if err != nil {
		return err
	}

	_, err = db.ExecContext(d.ctx, "INSERT INTO short_url (correlation_id, short_url) "+
		" VALUES ($1,$2)", uuid, code)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) Get(short string) (string, error) {
	db, err := sql.Open("pgx", d.DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	var uuid string

	row := db.QueryRowContext(d.ctx, "SELECT correlation_id FROM short_url WHERE short_url = $1", short)
	err = row.Scan(&uuid)
	if err != nil {
		log.Println(err)
	}

	var url URL

	row = db.QueryRowContext(d.ctx, "SELECT full_url FROM original_urls WHERE correlation_id = $1", uuid)
	err = row.Scan(&url.fullURL)
	if err != nil {
		log.Println(err)
	}

	return url.fullURL, nil
}

func (d *DB) GetExist(url string) (string, error) {
	db, err := sql.Open("pgx", d.DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	var id string

	row := db.QueryRowContext(d.ctx, "SELECT correlation_id FROM original_urls WHERE full_url = $1", url)
	err = row.Scan(&id)
	if err != nil {
		log.Println(err)
	}

	var short string

	row = db.QueryRowContext(d.ctx, "SELECT short_url FROM short_url WHERE correlation_id = $1", id)
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
