package clickhousetest

import (
	"context"
	"testing"
	"time"

	"github.com/rautaruukkipalich/urlsh/config"
	"github.com/rautaruukkipalich/urlsh/internal/store/clickhouse"
	"github.com/stretchr/testify/assert"
)

func Test_NewClickhouse(t *testing.T) {

	cfg := &config.DBConfig{
		URI:              "localhost:9000",
		Database:         "default",
		Username:         "default",
		Password:         "",
		MaxExecutionTime: 60,
		DialTimeout:      30 * time.Second,
	}
	db, err := clickhouse.New(context.TODO(), cfg)

	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func Test_CreateDatabase(t *testing.T) {
	ctx := context.TODO()
	db, teardown := testDB(t)

	defer teardown("test")

	err := db.Conn.Exec(ctx, "CREATE DATABASE IF NOT EXISTS test")
	if err != nil {
		t.Fatal(err)
	}

	err = db.Conn.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS test.urls (
				short String,
				long String
			) engine=Memory;
		`)

	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}

func Test_SetURLs(t *testing.T) {
	db, teardown := testDB(t)
	database, table := "test", "test"
	createTablesTestDB(t, db, database, table)
	defer teardown(database)

	storage := newStorage(db)
	ctx := context.TODO()

	var err error

	for _, urls := range testURLs {
		err := storage.setter.SetURLs(ctx, &urls)
		if err != nil {
			assert.NoError(t, err)
		}
	}

	assert.NoError(t, err)
}

func Test_GetShortURL(t *testing.T) {
	db, teardown := testDB(t)
	database, table := "test", "test"
	createTablesTestDB(t, db, database, table)
	defer teardown(database)

	storage := newStorage(db)
	ctx := context.TODO()

	for _, urls := range testURLs {
		_ = storage.setter.SetURLs(ctx, &urls)
	}

	for _, urls := range testURLs {
		urls.Short = ""
		ok, err := storage.getter.GetShortURL(ctx, &urls)
		assert.NoError(t, err)
		assert.NotEmpty(t, urls.Short)
		assert.True(t, ok)
	}
}

func Test_GetLongURL(t *testing.T) {
	db, teardown := testDB(t)
	database, table := "test", "test"
	createTablesTestDB(t, db, database, table)
	defer teardown(database)

	storage := newStorage(db)
	ctx := context.TODO()

	for _, urls := range testURLs {
		_ = storage.setter.SetURLs(ctx, &urls)
	}

	for _, urls := range testURLs {
		urls.Long = ""
		ok, err := storage.getter.GetLongURL(ctx, &urls)
		assert.NoError(t, err)
		assert.NotEmpty(t, urls.Long)
		assert.True(t, ok)
	}
}
