package clickhousetest

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rautaruukkipalich/urlsh/config"
	"github.com/rautaruukkipalich/urlsh/internal/store"
	"github.com/rautaruukkipalich/urlsh/internal/store/clickhouse"
)

func testDB(t *testing.T) (*clickhouse.DB, func(...string)) {
	t.Helper()

	ctx := context.TODO()

	cfg := &config.DBConfig{
		URI:              "localhost:9000",
		Database:         "default",
		Username:         "default",
		Password:         "",
		MaxExecutionTime: 60,
		DialTimeout:      30 * time.Second,
	}
	db, err := clickhouse.New(ctx, cfg)

	if err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			if err := db.Conn.Exec(
				ctx,
				fmt.Sprintf("DROP DATABASE %s", strings.Join(tables, ", ")),
			); err != nil {
				t.Fatal(err)
			}
		}
		db.Conn.Close()
	}
}

func createTablesTestDB(t *testing.T, db *clickhouse.DB, database, table string) {
	ctx := context.TODO()
	err := db.Conn.Exec(
		ctx,
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", database),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Conn.Exec(
		ctx,
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
				short String,
				long String
			) engine=Memory;
		`,
			database,
			table,
		),
	)

	if err != nil {
		t.Fatal(err)
	}
}

type testStore struct {
	setter store.StoreSetter
	getter store.StoreGetter
	stopper store.StoreStopper
}

func newStorage(db *clickhouse.DB) *testStore {
	return &testStore{
		setter: db,
		getter: db,
		stopper: db,
	}
}