package mocktest

import (
	"context"
	"testing"

	"github.com/rautaruukkipalich/urlsh/config"
	"github.com/rautaruukkipalich/urlsh/internal/store/mock"
	"github.com/stretchr/testify/assert"
)

func Test_NewNock(t *testing.T) {
	ctx := context.Background()
	cfg := config.DBConfig{}
	db, _ := mock.New(ctx, &cfg)

	assert.NotNil(t, db)
}

func Test_SetURLs(t *testing.T) {
	db, teardown := testDB(t)
	defer teardown()

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
	defer teardown()

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
	defer teardown()

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
