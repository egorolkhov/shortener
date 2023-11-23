package app

import (
	"github.com/egorolkhov/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_ShortURL(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		storageSize int
		method      string
	}
	tests := []struct {
		name    string
		fields  *storage.Data
		want    want
		request string
	}{
		{
			name:   "simple test",
			fields: &storage.Data{Urls: map[string]string{}},
			want: want{
				contentType: "text/plain",
				statusCode:  201,
				storageSize: 1,
				method:      http.MethodPost,
			},
			request: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{tt.fields, "", "", "", 0}

			request := httptest.NewRequest(tt.want.method, "/", nil)

			w := httptest.NewRecorder()

			a.ShortURL(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.statusCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"), "Тип ответа не совпадает")
			assert.Equal(t, tt.want.storageSize, len(a.Storage.(*storage.Data).Urls))
			defer res.Body.Close()
		})
	}
}
