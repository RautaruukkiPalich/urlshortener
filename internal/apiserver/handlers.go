package apiserver

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rautaruukkipalich/urlsh/internal/model"
	logger "github.com/rautaruukkipalich/urlsh/pkg/log"
)

// @Summary		Push url, get short url
// @Description	push url, get short url
// @Tags			urls
// @Accept			json
// @Produce		json
// @Param			input	body		model.URL	true	"url from"
// @Success		200		{object}	model.URLs
// @Failure		400,404	{object}	errorResponse
// @Success		500		{object}	errorResponse
// @Success		default	{object}	errorResponse
// @Router			/shorten [post]
func (a *APIServer) PushLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "internal.apiserver.PushLink"
		ctx := r.Context()

		ctx = logger.AddTrace(ctx, slog.Any("op", op))
		ctx = logger.AddGroup(ctx, slog.Group(
			"req",
			slog.String("request uuid", uuid.New().String()),
			slog.String("remote addr", r.RemoteAddr),
			slog.String("request addr", r.RequestURI),
		))

		defer r.Body.Close()

		var urls model.URLs

		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			logger.LoggerFromContext(ctx).Error("decode json error", slog.String("err", err.Error()))
			a.error(ctx, w, r, ErrResp500)
			return
		}

		ctx = logger.AddAttr(ctx, "longurl", urls.Long)

		if err := a.vaidateUrl(urls.Long); err != nil {
			logger.LoggerFromContext(ctx).Error("invalid url", slog.String("err", err.Error()))
			a.error(ctx, w, r, errorResponse{http.StatusBadRequest, ErrInvalidURL.Error()})
			return
		}

		ok, err := a.GetShortURL(ctx, &urls)
		if err != nil {
			if err != sql.ErrNoRows {
				a.error(ctx, w, r, ErrResp500)
				return
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

				ok, err := a.GetLongURL(ctx, &urls)
				if err != nil {
					logger.LoggerFromContext(ctx).Debug("err while get long url", slog.String("err", err.Error()))
					a.error(ctx, w, r, ErrResp500)
					return
				}

				if ok {
					continue
				}

				if err := a.SetURLs(ctx, &urls); err != nil {
					logger.LoggerFromContext(ctx).Debug("error while set urls error", slog.String("err", err.Error()))
					continue
				}
				logger.LoggerFromContext(ctx).Debug("short url generated")
				break
			}
		}

		a.JSONrespond(ctx, w, r, http.StatusOK, urls)
	}
}

// @Summary		Get long url from short url
// @Description	get long url from short url
// @Tags			urls
// @Accept			json
// @Produce		json
// @Param			url		path		string	true	"url"
// @Success		200		{object}	model.URLs
// @Failure		400,404	{object}	errorResponse
// @Success		500		{object}	errorResponse
// @Success		default	{object}	errorResponse
// @Router			/{url} [get]
func (a *APIServer) LinkFromShortUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "internal.apiserver.LinkFromShortUrl"
		ctx := logger.AddTrace(r.Context(), slog.Any("op", op))

		ctx = logger.AddGroup(ctx, slog.Group(
			"req",
			slog.String("request uuid", uuid.New().String()),
			slog.String("remote addr", r.RemoteAddr),
			slog.String("request addr", r.RequestURI),
		))

		var urls model.URLs

		short := strings.TrimLeft(r.RequestURI, "/")
		urls.Short = short

		ok, err := a.GetLongURL(ctx, &urls)
		if err != nil || !ok {
			if err == nil {
				err = ErrNotFound
			} 
			logger.LoggerFromContext(ctx).Debug("error while get cached long url", slog.Any("error", err.Error()))
			a.error(ctx, w, r, errorResponse{Code: http.StatusBadRequest, Error: ErrInvalidURL.Error()})
			return
		}

		a.JSONrespond(ctx, w, r, http.StatusOK, urls)
	}
}

// func (a *APIServer) RedirectShortUrl() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		const op = "internal.apiserver.RedirectShortUrl"
// 		log := logger.LoggerFromContext(a.ctx).With(
// 			slog.String("op", op),
// 			slog.String("remote addr", r.RemoteAddr),
// 			slog.String("request addr", r.RequestURI),
// 		)
// 		// log := a.log.With(
// 		// 	slog.String("op", op),
// 		// 	slog.String("remote addr", r.RemoteAddr),
// 		// 	slog.String("request addr", r.RequestURI),
// 		// )

// 		var urls model.URLs

// 		defer r.Body.Close()

// 		pattern := strings.TrimLeft(r.RequestURI, "/")
// 		urls.Short = pattern

// 		ok, err := a.GetLongURL(a.ctx, &urls)
// 		if err != nil || !ok {
// 			log.Error("some error", slog.String("err", err.Error()))
// 			a.error(w, r, errorResponse{Code: http.StatusBadRequest, Error: ErrInvalidURL.Error()})
// 			return
// 		}

// 		a.redirect(w, r, urls.Long)
// 	}
// }
