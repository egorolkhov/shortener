package app

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func (a *App) PSQLconnection(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		a.DatabaseDSN.Host, a.DatabaseDSN.Port, a.DatabaseDSN.User, a.DatabaseDSN.Password, a.DatabaseDSN.DBname, a.DatabaseDSN.SSLmode)

	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		return
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
