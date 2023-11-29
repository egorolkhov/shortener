package app

import (
	"encoding/json"
	"fmt"
	"github.com/egorolkhov/shortener/internal/middleware"
	"net/http"
)

func (a *App) UserAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie := w.Header().Get("Authorization")
	userID := middleware.GetUserID(cookie, "1234")
	fmt.Println("TOKEN", cookie)
	fmt.Println("USERID", userID)
	urls := a.Storage.GetUserURLS(userID)
	fmt.Println("URLS", urls)

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println(string(result))
	w.Write(result)
}
