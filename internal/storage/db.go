package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type URL struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}

type DB struct {
	ctx         context.Context
	DatabaseDSN string
}

func NewDB(DatabaseDSN string) *DB {
	return &DB{context.Background(), DatabaseDSN}
}

func (d *DB) Add(userID, code, url string) error {
	db, err := sql.Open("pgx", d.DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	uuid := generateCorrelationID()

	result, err := db.ExecContext(d.ctx, "INSERT INTO original_urls (user_id, correlation_id, full_url)"+
		" VALUES ($1,$2,$3) ON CONFLICT (full_url) DO NOTHING;", userID, uuid, url)

	if i, _ := result.RowsAffected(); i == 0 {
		return ErrURLAlreadyExist
	}
	if err != nil {
		return err
	}

	_, err = db.ExecContext(d.ctx, "INSERT INTO short_urls (user_id, correlation_id, short_url) "+
		" VALUES ($1,$2,$3)", userID, uuid, code)
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

	row := db.QueryRowContext(d.ctx, "SELECT correlation_id FROM short_urls WHERE short_url = $1", short)
	err = row.Scan(&uuid)
	if err != nil {
		log.Println("HERE1", err)
	}

	var url URL

	row = db.QueryRowContext(d.ctx, "SELECT full_url FROM original_urls WHERE correlation_id = $1", uuid)
	err = row.Scan(&url.FullURL)
	if err != nil {
		log.Println("HERE2", err)
	}

	return url.FullURL, nil
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

	row = db.QueryRowContext(d.ctx, "SELECT short_url FROM short_urls WHERE correlation_id = $1", id)
	err = row.Scan(&short)
	if err != nil {
		log.Println(err)
	}

	return short, nil
}

func AddBatch(ctx context.Context, DatabaseDSN string, userID string, codes []string, jsons []RequestBatch) error {
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

	stmtShort, err := tx.PrepareContext(ctx, "INSERT INTO short_urls (user_id, correlation_id, short_url)"+
		" VALUES ($1,$2,$3)")
	if err != nil {
		return err
	}
	defer stmtShort.Close()

	stmtOriginal, err := tx.PrepareContext(ctx, "INSERT INTO original_urls (user_id, correlation_id, full_url) "+
		"VALUES($1,$2,$3)")
	if err != nil {
		return err
	}
	defer stmtOriginal.Close()

	for i, line := range jsons {
		_, err = stmtShort.ExecContext(ctx, userID, line.CorrelationID, codes[i])
		if err != nil {
			return err
		}
		_, err = stmtOriginal.ExecContext(ctx, userID, line.CorrelationID, line.OriginalURL)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d *DB) GetUserURLS(userID string) []URL {
	db, err := sql.Open("pgx", d.DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	var urls []URL
	var full []string
	var short []string

	rows, err := db.QueryContext(d.ctx, "SELECT full_url FROM original_urls WHERE user_id = $1", userID)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var fullURL string
		if err = rows.Scan(&fullURL); err != nil {
			log.Fatal(err)
		}
		full = append(full, fullURL)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}

	rows, err = db.QueryContext(d.ctx, "SELECT short_url FROM short_urls WHERE user_id = $1", userID)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var shortURL string
		if err = rows.Scan(&shortURL); err != nil {
			log.Fatal(err)
		}
		short = append(short, shortURL)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}

	if len(full) != 0 {
		for i := 0; i < len(full); i++ {
			urls = append(urls, URL{"http://localhost:8080/" + short[i], full[i]})
		}
	}
	return urls
}

func IsDeleted(ctx context.Context, DatabaseDSN string, short string) bool {
	db, err := sql.Open("pgx", DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	var deleted bool

	row := db.QueryRowContext(ctx, "SELECT is_deleted FROM short_urls WHERE short_url = $1", short)
	err = row.Scan(&deleted)
	if err != nil {
		log.Println(err)
	}
	return deleted
}

func DeleteURL(ctx context.Context, DatabaseDSN string, userID string, codes []string) error {
	db, err := sql.Open("pgx", DatabaseDSN)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	var uuids []string

	rows, err := db.QueryContext(ctx, "SELECT correlation_id FROM short_urls WHERE user_id = $1 and short_url = ANY($2)",
		userID, codes)
	if err != nil {
		log.Println(err)
		return err
	}
	for rows.Next() {
		var uuid string
		if err = rows.Scan(&uuid); err != nil {
			log.Fatal(err)
		}
		uuids = append(uuids, uuid)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	tx, err := db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	updt, err := tx.PrepareContext(ctx, "UPDATE short_urls SET is_deleted = true WHERE correlation_id = $1")
	if err != nil {
		return err
	}
	defer updt.Close()

	fmt.Println(uuids)
	for _, uuid := range uuids {
		_, err = updt.Exec(uuid)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
