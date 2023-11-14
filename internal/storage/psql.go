package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/egorolkhov/shortener/internal/config"
	"log"
)

type URL struct {
	fullURL  string
	shortURL string
}

func CreateTable(DatabaseDSN config.PGXaddress) error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DatabaseDSN.Host, DatabaseDSN.Port, DatabaseDSN.User, DatabaseDSN.Password, DatabaseDSN.DBname, DatabaseDSN.SSLmode)

	db, err := sql.Open("pgx", psqlInfo)
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

func AddDB(ctx context.Context, DatabaseDSN config.PGXaddress, code, url string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DatabaseDSN.Host, DatabaseDSN.Port, DatabaseDSN.User, DatabaseDSN.Password, DatabaseDSN.DBname, DatabaseDSN.SSLmode)

	db, err := sql.Open("pgx", psqlInfo)
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

func GetDB(ctx context.Context, DatabaseDSN config.PGXaddress, short string) (string, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DatabaseDSN.Host, DatabaseDSN.Port, DatabaseDSN.User, DatabaseDSN.Password, DatabaseDSN.DBname, DatabaseDSN.SSLmode)

	db, err := sql.Open("pgx", psqlInfo)
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
