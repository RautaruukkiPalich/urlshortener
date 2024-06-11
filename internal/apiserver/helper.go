package apiserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/rautaruukkipalich/urlsh/internal/model"
	logger "github.com/rautaruukkipalich/urlsh/pkg/log"
)

const (
	Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@-_"
	URLlength = 8
)

func (a *APIServer) error(ctx context.Context, w http.ResponseWriter, r *http.Request, code int, err error) {
	a.JSONrespond(ctx, w, r, code, errorResponse{Code: code, Error: err.Error()})
}

func (a *APIServer) JSONrespond(ctx context.Context, w http.ResponseWriter, r *http.Request, code int, data any) {
	const op = "internal.apiserver.JSONrespond"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))

	w.WriteHeader(code)

	if data != nil {
		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.LoggerFromContext(ctx).Error("error encode json", slog.Any("error", err))
		}
	}
}

func (a *APIServer) generateShortUrl() string {
	b := make([]byte, URLlength)
	for i := range b {
		b[i] = Letters[rand.Intn(len(Letters))]
	}
	
	return string(b)
}

func (a *APIServer) vaidateUrl(url string) error {
	if len(url) == 0 {
		return ErrInvalidURL
	}

	pattern := `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`

	ok, err := regexp.Match(pattern, []byte(url))
	if err != nil || !ok {
		return ErrInvalidURL
	}

	return nil
}


func (a *APIServer) redirect(ctx context.Context, w http.ResponseWriter, r *http.Request, code int, urls *model.URLs) {
	logger.LoggerFromContext(ctx).Info("reditected")
	http.Redirect(w, r, urls.Long, http.StatusPermanentRedirect)
}


func (a *APIServer) GetShortURL(ctx context.Context, urls *model.URLs) (ok bool, err error){
	const op = "internal.apiserver.GetShortURL"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))

	ok, err = a.CachedGetLongURL(ctx, urls)

	if !ok && err == nil {
		logger.LoggerFromContext(ctx).Debug("long url from short url not found")
		return false, ErrNotFound
	}
	if err != nil {
		logger.LoggerFromContext(ctx).Debug("error while get cached long url", slog.Any("error", err.Error()))
		return false, ErrInternalServerError
	}

	return true, nil
}

func (a *APIServer) PushURL(ctx context.Context, urls *model.URLs) (int, error){
	const op = "internal.apiserver.PushURL"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))

	ok, err := a.CachedGetShortURL(ctx, urls)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LoggerFromContext(ctx).Debug("err while get short url from cache", slog.String("err", err.Error()))
			return http.StatusInternalServerError, ErrInternalServerError
		}
	}

	if ok {
		ctx = logger.AddAttr(ctx, "shorturl", urls.Short)
	}

	if urls.Short == EmptyString {
		for {
			short := a.generateShortUrl()
			urls.Short = short

			ctx = logger.AddAttr(ctx, "shorturl", urls.Short)

			ok, err = a.CachedGetLongURL(ctx, urls)
			if err != nil {
				logger.LoggerFromContext(ctx).Debug("err while get long url", slog.String("err", err.Error()))
				return http.StatusInternalServerError, ErrInternalServerError
			}

			if ok {
				continue
			}

			if err := a.CachedSetURLs(ctx, urls); err != nil {
				logger.LoggerFromContext(ctx).Debug("error while set urls error", slog.String("err", err.Error()))
				continue
			}
			logger.LoggerFromContext(ctx).Debug("short url generated")
			break
		}
	}

	return http.StatusOK, nil
}