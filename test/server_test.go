package testserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rautaruukkipalich/urlsh/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_PushLink(t *testing.T) {
	ctx := context.Background()
	srv := testServer(t, ctx)

	testCases := []struct {
		name         string
		payload      any
		expectedCode int
	}{
		{
			name: "invalid url 1",
			payload: map[string]string{
				"long": "123",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid url 2",
			payload: map[string]string{
				"long": "tcp:127.0.0.1:443",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid url 3",
			payload: map[string]string{
				"long": "t.me/123",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid payload",
			payload: map[string]string{
				"longg": "http://example.com",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "valid payload",
			payload: map[string]string{
				"long": "http://example.com",
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			var json_data bytes.Buffer
			_ = json.NewEncoder(&json_data).Encode(tc.payload)

			req := httptest.NewRequest(http.MethodPost, "/shorten", &json_data)
			srv.PushURLHandler().ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)
		})
	}
}

func Test_LinkFromShortUrl(t *testing.T) {
	ctx := context.Background()
	srv := testServer(t, ctx)

	testCases := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{
			name:         "invalid url",
			path:         "/21321",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid url 2",
			path:         "/21321/eqwewqe",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			srv.GetShortURLHandler().ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)
		})
	}
}

func Test_GetShortUrls(t *testing.T) {
	ctx := context.Background()
	srv := testServer(t, ctx)
	fillTestDB(t, srv)

	for _, tc := range URLsTestCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tc.urls.Short), nil)
			srv.GetShortURLHandler().ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)

			var urls model.URLs
			if err := json.NewDecoder(rr.Body).Decode(&urls); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.urls.Long, urls.Long)
		})
	}

}
