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

func CreateTable(DatabaseDSN string) error {
	db, err := sql.Open("pgx", DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	createTableSQL := `
        CREATE TABLE IF NOT EXISTS shortener (
            full_url TEXT,
            short_url TEXT
        )
    `

	_, err = db.Exec(createTableSQL)

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

	_, err = db.Exec("INSERT INTO shortener (short_url, full_url)"+
		" VALUES ($1,$2)", code, url)

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

	var url URL

	row := db.QueryRowContext(ctx, "SELECT full_url FROM shortener WHERE short_url = $1", short)

	err = row.Scan(&url.fullURL)
	if err != nil {
		log.Println(err)
	}

	return url.fullURL, nil
}
