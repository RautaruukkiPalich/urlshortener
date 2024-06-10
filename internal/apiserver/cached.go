package apiserver

import (
	"context"
	"log/slog"

	"github.com/rautaruukkipalich/urlsh/internal/cache"
	"github.com/rautaruukkipalich/urlsh/internal/model"
	logger "github.com/rautaruukkipalich/urlsh/pkg/log"
)

func (a *APIServer) GetLongURL(ctx context.Context, urls *model.URLs) (bool, error) {
	const op = "internal.apiserver.GetLongURL"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))

	ok, err := a.cache.get.GetLongURL(ctx, urls)
	if err != nil {
		if err != cache.ErrNotFound{
			return false, err
		}
	}
	if ok {
		ctx = logger.AddAttr(ctx, "longurl", urls.Long)
		logger.LoggerFromContext(ctx).Info("longurl form cache")
		return true, nil
	}

	ok, err = a.store.get.GetLongURL(ctx, urls)
	if err != nil || !ok{
		return false, err
	}

	if err := a.cache.set.SetURLs(ctx, urls); err != nil {
		return false, err
	}

	ctx = logger.AddAttr(ctx, "longurl", urls.Long)
	logger.LoggerFromContext(ctx).Info("cached longurl")

	return true, nil
}

func (a *APIServer) GetShortURL(ctx context.Context, urls *model.URLs) (bool, error) {
	const op = "internal.apiserver.GetShortURL"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))

	ok, err := a.cache.get.GetShortURL(ctx, urls)
	if err != nil {
		if err != cache.ErrNotFound{
			return false, err
		}
	}
	if !ok {
		ok, err = a.store.get.GetShortURL(ctx, urls)
		if err != nil || !ok{
			return false, err
		}

		if err := a.cache.set.SetURLs(ctx, urls); err != nil {
			return false, err
		}
		logger.LoggerFromContext(ctx).Info("urls cached")
	}

	return true, nil
}

func (a *APIServer) SetURLs(ctx context.Context, urls *model.URLs) (err error) {
	const op = "internal.apiserver.SetURLs"
	ctx = logger.AddTrace(ctx, slog.String("op", op))

	if err := a.store.set.SetURLs(ctx, urls); err != nil {
		return err
	}

	if err := a.cache.set.SetURLs(ctx, urls); err != nil {
		return err
	}

	logger.LoggerFromContext(ctx).Info("urls cached")

	return nil
}
