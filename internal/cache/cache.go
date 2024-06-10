package cache

import (
	"context"

	"github.com/rautaruukkipalich/urlsh/internal/model"
)

type CacheGetter interface {
	GetLongURL(ctx context.Context, urls *model.URLs) (bool, error)
	GetShortURL(ctx context.Context, urls *model.URLs) (bool, error)
}

type CacheSetter interface {
	SetURLs(ctx context.Context, urls *model.URLs) error
}

type CacheStopper interface {
	Stop(context.Context)
}
