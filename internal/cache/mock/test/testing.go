package mockcachetest

import (
	"context"
	"testing"

	"github.com/rautaruukkipalich/urlsh/internal/cache"
	mockcache "github.com/rautaruukkipalich/urlsh/internal/cache/mock"
)

type MockCache struct {
	getter cache.CacheGetter
	setter cache.CacheSetter
	stoper cache.CacheStopper
}

func testMockCache(t *testing.T) (*MockCache, error) {
	t.Helper()

	client, err := mockcache.New(context.TODO(), &cfg)
	if err != nil {
		t.Fatal(err)
	}

	
	return &MockCache{
		getter: client,
		setter: client,
		stoper: client,
	}, nil
}