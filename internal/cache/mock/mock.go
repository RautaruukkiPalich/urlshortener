package mockcache

import (
	"context"

	"github.com/rautaruukkipalich/urlsh/config"
	"github.com/rautaruukkipalich/urlsh/internal/model"
)

type Cache struct {
	store map[string]string
}

func New(ctx context.Context, cfg *config.CacheConfig) (*Cache, error) {
	return &Cache{store: make(map[string]string)}, nil
}

func (c *Cache) GetLongURL(ctx context.Context, urls *model.URLs) (bool, error) {

	longurl, ok := c.store[urls.Short]
	if ok {
		urls.Long = longurl
	}
	return ok, nil
}

func (c *Cache) GetShortURL(ctx context.Context, urls *model.URLs) (bool, error) {
	
	shorturl, ok := c.store[urls.Long]
	if ok {
		urls.Short = shorturl
	}
	return ok, nil
}


func (c *Cache) SetURLs(ctx context.Context, urls *model.URLs) error {
	
	c.store[urls.Short] = urls.Long
	c.store[urls.Long] = urls.Short
	return nil
}

func (c *Cache) Stop(ctx context.Context) {
}