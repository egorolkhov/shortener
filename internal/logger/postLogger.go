package logger

import (
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type responseData struct {
	StatusCode int
	Size       int
}

type LogResponseWriter struct {
	http.ResponseWriter
	Log *responseData
}

func (L *LogResponseWriter) Write(b []byte) (int, error) {
	size, err := L.ResponseWriter.Write(b)
	if err != nil {
		return 0, err
	}
	L.Log.Size = size
	return size, nil
}

func (L *LogResponseWriter) WriteHeader(statusCode int) {
	L.ResponseWriter.WriteHeader(statusCode)
	L.Log.StatusCode = statusCode
}

func PostLogger(h http.HandlerFunc) http.HandlerFunc {
	Foo := func(w http.ResponseWriter, r *http.Request) {
		Data := &responseData{0, 0}

		lw := LogResponseWriter{w, Data}
		h(&lw, r)

		Log.Info(
			"INFO",
			zap.String("status code", strconv.Itoa(Data.StatusCode)),
			zap.String("size", strconv.Itoa(Data.Size)),
		)
	}
	return Foo
}
