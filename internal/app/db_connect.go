package app

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func (a *App) PSQLconnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println(a.DatabaseDSN)

	db, err := sql.Open("pgx", a.DatabaseDSN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
