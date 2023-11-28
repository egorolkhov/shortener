package app

import (
	"encoding/json"
	"fmt"
	"github.com/egorolkhov/shortener/internal/middleware"
	"net/http"
)

func (a *App) UserAPI(w http.ResponseWriter, r *http.Request) {
	cookie := w.Header().Get("Authorization")
	userID := middleware.GetUserID(cookie, "1234")
	if userID == "error" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	urls := a.Storage.GetUserURLS(userID)

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return

	}

	fmt.Println(urls)

	result, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(string(result))
	w.Write(result)
}

type URL struct {
	FullURL string
}

type User struct {
	Name string
}
