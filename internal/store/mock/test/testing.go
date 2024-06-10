package mocktest

import (
	"context"
	"testing"

	"github.com/rautaruukkipalich/urlsh/config"
	"github.com/rautaruukkipalich/urlsh/internal/store"
	"github.com/rautaruukkipalich/urlsh/internal/store/mock"
)

func testDB(t *testing.T) (*mock.DB, func(...string)) {
	t.Helper()
	ctx := context.Background()
	cfg := config.DBConfig{}
	db, _ := mock.New(ctx, &cfg)
	return db, func(...string) {}
}

type testStore struct {
	setter  store.StoreSetter
	getter  store.StoreGetter
	stopper store.StoreStopper
}

func newStorage(db *mock.DB) *testStore {
	return &testStore{
		setter:  db,
		getter:  db,
		stopper: db,
	}
}