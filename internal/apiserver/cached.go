package apiserver

import (
	"context"
	"log/slog"

	"github.com/rautaruukkipalich/urlsh/internal/cache"
	"github.com/rautaruukkipalich/urlsh/internal/model"
	"github.com/rautaruukkipalich/urlsh/internal/store"
	logger "github.com/rautaruukkipalich/urlsh/pkg/log"
)

var (
	FailGetLongUrlFromCache = "fail get longurl form cache"
	FailGetLongUrlFromStore = "fail get longurl form store"
	
	LongUrlFromCache        = "long url form cache"
	ShortUrlFromCache       = "short url form cache"
	
	NotFoundInStore         = "not found in store"
	NotFoundInCache         = "not found in cache"

	FailSetURLsToCache      = "fail set urls to cache"
	FailSetURLsToStore      = "fail set urls to store"
	
	LongUrlCached           = "longurl cached"
	URLsCached              = "urls cached"
)

func (a *APIServer) CachedGetLongURL(ctx context.Context, urls *model.URLs) (bool, error) {
	const op = "internal.apiserver.CachedGetLongURL"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))

	ok, err := a.cache.get.GetLongURL(ctx, urls)

	if err != nil {
		switch err {
		case cache.ErrNotFound:
		default:
			logger.LoggerFromContext(ctx).Debug(FailGetLongUrlFromCache)
			return false, err
		}
	}

	if ok {
		ctx = logger.AddAttr(ctx, "longurl", urls.Long)
		logger.LoggerFromContext(ctx).Debug(LongUrlFromCache)
		return ok, nil
	}

	ok, err = a.store.get.GetLongURL(ctx, urls)
	if err != nil {
		switch err {
		case store.ErrNotFound:
		default:
			logger.LoggerFromContext(ctx).Debug(FailGetLongUrlFromStore)
			return false, ErrInternalServerError
		}
	}

	if !ok {
		logger.LoggerFromContext(ctx).Debug(NotFoundInStore)
		return ok, nil
	}

	ctx = logger.AddAttr(ctx, "longurl", urls.Long)

	if err := a.cache.set.SetURLs(ctx, urls); err != nil {
		logger.LoggerFromContext(ctx).Debug(FailSetURLsToCache)
		return false, err
	}
	logger.LoggerFromContext(ctx).Debug(URLsCached)

	return true, nil
}

func (a *APIServer) CachedGetShortURL(ctx context.Context, urls *model.URLs) (bool, error) {
	const op = "internal.apiserver.CachedGetShortURL"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))

	ok, err := a.cache.get.GetShortURL(ctx, urls)
	if err != nil {
		switch err {
		case cache.ErrNotFound:
		default:
			return false, err
		}
	}

	if !ok {
		ok, err = a.store.get.GetShortURL(ctx, urls)
		if err != nil || !ok {
			return false, err
		}

		if err := a.cache.set.SetURLs(ctx, urls); err != nil {
			logger.LoggerFromContext(ctx).Debug(FailSetURLsToCache)
			return false, err
		}
		logger.LoggerFromContext(ctx).Debug(URLsCached)
	}

	return true, nil
}

func (a *APIServer) CachedSetURLs(ctx context.Context, urls *model.URLs) (err error) {
	const op = "internal.apiserver.CachedSetURLs"
	ctx = logger.AddTrace(ctx, slog.String("op", op))

	if err := a.store.set.SetURLs(ctx, urls); err != nil {
		logger.LoggerFromContext(ctx).Debug(FailSetURLsToStore)
		return err
	}

	if err := a.cache.set.SetURLs(ctx, urls); err != nil {
		logger.LoggerFromContext(ctx).Debug(FailSetURLsToCache)
		return err
	}

	logger.LoggerFromContext(ctx).Debug(URLsCached)

	return nil
}
