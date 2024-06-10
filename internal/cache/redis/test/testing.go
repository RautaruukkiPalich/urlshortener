package redistest

import (
	"context"
	"testing"

	"github.com/rautaruukkipalich/urlsh/internal/cache"
	rediscache "github.com/rautaruukkipalich/urlsh/internal/cache/redis"
)

type RedisCache struct {
	getter cache.CacheGetter
	setter cache.CacheSetter
	stoper cache.CacheStopper
}

func testRedisCache(t *testing.T) (*RedisCache, error) {
	t.Helper()

	client, err := rediscache.New(context.TODO(), &cfg)
	if err != nil {
		t.Fatal(err)
	}
	
	return &RedisCache{
		getter: client,
		setter: client,
		stoper: client,
	}, nil
}