package store

import (
	"context"

	"github.com/rautaruukkipalich/urlsh/internal/model"
)

type StoreStopper interface {
	Stop(context.Context)
}

type StoreGetter interface {
	GetShortURL(context.Context, *model.URLs) (bool, error)
	GetLongURL(context.Context, *model.URLs) (bool, error)
}

type StoreSetter interface {
	SetURLs(context.Context, *model.URLs) error
}
