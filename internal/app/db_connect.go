package app

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func (a *App) PSQLconnection(w http.ResponseWriter, r *http.Request) {
	//psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	//	a.DatabaseDSN.Host, a.DatabaseDSN.Port, a.DatabaseDSN.User, a.DatabaseDSN.Password, a.DatabaseDSN.DBname, a.DatabaseDSN.SSLmode)
	psqlInfo := fmt.Sprintf("host=%s, user=%s password=%s dbname=%s ",
		a.DatabaseDSN.Host, a.DatabaseDSN.User, a.DatabaseDSN.Password, a.DatabaseDSN.DBname)

	db, err := sql.Open("pgx", psqlInfo)
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
