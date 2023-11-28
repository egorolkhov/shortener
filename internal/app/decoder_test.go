package app

import (
	"github.com/egorolkhov/shortener/internal/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_DecodeURL(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		url         string
		method      string
	}
	tests := []struct {
		name    string
		fields  *storage.Data
		want    want
		request string
	}{
		{
			name: "simple test",
			fields: &storage.Data{
				Urls:  map[string]string{"https://example.com": "cdFCblAL"},
				Codes: map[string]string{"cdFCblAL": "https://example.com"},
				Users: map[string][]storage.URL{"cdFCblAL": []storage.URL{{"short", "fillURL"}}}},
			want: want{
				contentType: "text/plain",
				statusCode:  307,
				url:         "https://example.com",
				method:      http.MethodGet,
			},
			request: "cdFCblAL",
		},
		{
			name: "wrong code",
			fields: &storage.Data{
				Urls:  map[string]string{"https://example.com": "cdFCblAL"},
				Codes: map[string]string{"cdFCblAL": "https://example.com"},
				Users: map[string][]storage.URL{}},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				method:      http.MethodGet,
			},
			request: "cdFCblA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{tt.fields, "", "", "", 0}

			request := httptest.NewRequest(tt.want.method, "/"+tt.request, nil)

			w := httptest.NewRecorder()

			vars := map[string]string{
				"id": tt.request,
			}

			request = mux.SetURLVars(request, vars)

			a.DecodeURL(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"), "Тип ответа не совпадает")
			assert.Equal(t, tt.want.url, res.Header.Get("Location"), "Ссылки не совпадают")
			defer res.Body.Close()
		})
	}
}
