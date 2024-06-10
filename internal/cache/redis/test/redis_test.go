package redistest

import (
	"context"
	"testing"

	rediscache "github.com/rautaruukkipalich/urlsh/internal/cache/redis"
	"github.com/stretchr/testify/assert"
)

func Test_NewRedis(t *testing.T) {

	cache, err := rediscache.New(context.TODO(), &cfg)

	assert.NoError(t, err)
	assert.NotNil(t, cache)
}

func Test_SetURLs(t *testing.T) {
	cache, err := testRedisCache(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()


	for _, urls := range testURLs {
		err := cache.setter.SetURLs(ctx, &urls)
		if err != nil {
			assert.NoError(t, err)
		}
	}

	assert.NoError(t, err)
}

func Test_GetShortURL(t *testing.T) {
	cache, err := testRedisCache(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	for _, urls := range testURLs {
		_ = cache.setter.SetURLs(ctx, &urls)
	}

	for _, urls := range testURLs {
		urls.Short = ""
		ok, err := cache.getter.GetShortURL(ctx, &urls)
		assert.NoError(t, err)
		assert.NotEmpty(t, urls.Short)
		assert.True(t, ok)
	}
}

func Test_GetLongURL(t *testing.T) {
	cache, err := testRedisCache(t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	for _, urls := range testURLs {
		_ = cache.setter.SetURLs(ctx, &urls)
	}

	for _, urls := range testURLs {
		urls.Long = ""
		ok, err := cache.getter.GetLongURL(ctx, &urls)
		assert.NoError(t, err)
		assert.NotEmpty(t, urls.Long)
		assert.True(t, ok)
	}
}

func Test_Stop(t *testing.T) {
	cache, err := testRedisCache(t)
	if err != nil {
		t.Fatal(err)
	}
	cache.stoper.Stop(context.Background())
}