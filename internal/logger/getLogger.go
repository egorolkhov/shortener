package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func GetLogger(h http.HandlerFunc) http.HandlerFunc {
	Foo := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			Time := time.Now()
			h(w, r)
			Duration := time.Since(Time)
			Log.Info(
				"INFO",
				zap.String("method", r.Method),
				zap.String("time", Duration.String()),
				zap.String("URI", r.RequestURI),
			)
		}
	}
	return Foo
}
