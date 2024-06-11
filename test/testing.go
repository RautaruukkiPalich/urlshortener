package testserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rautaruukkipalich/urlsh/config"
	"github.com/rautaruukkipalich/urlsh/internal/apiserver"
	mockcache "github.com/rautaruukkipalich/urlsh/internal/cache/mock"
	"github.com/rautaruukkipalich/urlsh/internal/model"
	"github.com/rautaruukkipalich/urlsh/internal/store/mock"
	"github.com/rautaruukkipalich/urlsh/pkg/loggerloader"
	logger "github.com/rautaruukkipalich/urlsh/pkg/log"
)

func testServer(t *testing.T, ctx context.Context) *apiserver.APIServer {
	t.Helper()

	log := loggerloader.MustRunLogger("prod")
	logger.CreateGlobalLogger(log)

	var (
		cfg = &config.SRVConfig{
			Addr:    "localhost:8080",
			Timeout: 10 * time.Second,
		}
	)

	store, err := mock.New(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	cache, err := mockcache.New(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	srv := apiserver.New(ctx, cfg, store, store, store, cache, cache, cache)

	return srv
}

func fillTestDB(t *testing.T, srv *apiserver.APIServer) {
	for idx, tc := range URLsTestCases {
		rr := httptest.NewRecorder()

		var json_data bytes.Buffer
		_ = json.NewEncoder(&json_data).Encode(tc.urls)

		req := httptest.NewRequest(http.MethodPost, "/shorten", &json_data)
		srv.PushURLHandler().ServeHTTP(rr, req)

		var urls model.URLs 
		if err := json.NewDecoder(rr.Body).Decode(&urls); err != nil {
			t.Fatal(err)
			return
		}
		tc.urls = urls
		URLsTestCases[idx] = tc
	}
}
