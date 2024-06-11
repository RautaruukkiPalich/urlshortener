package apiserver

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

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
func (a *APIServer) PushURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var code int = http.StatusOK

		start := time.Now()
		defer func(){
			observeRequest(time.Since(start), code)
		}()

		const op = "internal.apiserver.PushURLHandler"
		ctx := r.Context()

		ctx = logger.AddTrace(ctx, slog.Any("op", op))
		ctx = logger.AddGroup(ctx, slog.Group(
			"req",
			slog.String("request uuid", uuid.New().String()),
			slog.String("remote addr", r.RemoteAddr),
			slog.String("request addr", r.RequestURI),
		))

		var urls model.URLs

		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			logger.LoggerFromContext(ctx).Error("decode json error", slog.String("err", err.Error()))
			code = http.StatusInternalServerError
			a.error(ctx, w, r, code, ErrInternalServerError)
			return
		}

		ctx = logger.AddAttr(ctx, "longurl", urls.Long)

		if err := a.vaidateUrl(urls.Long); err != nil {
			logger.LoggerFromContext(ctx).Error("invalid url", slog.String("err", err.Error()))
			code = http.StatusBadRequest
			a.error(ctx, w, r, code, err)
			return
		}

		code, err := a.PushURL(ctx, &urls)
		if err != nil{
			a.error(ctx, w, r, code, err)
		} 

		a.JSONrespond(ctx, w, r, code, urls)
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
func (a *APIServer) GetShortURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var code int = http.StatusOK

		start := time.Now()
		defer func(){
			observeRequest(time.Since(start), code)
		}()

		const op = "internal.apiserver.GetShortURLHandler"
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

		ok, err := a.GetShortURL(ctx, &urls)
		if err != nil {
			code = http.StatusInternalServerError
			a.error(ctx, w, r, code, err)
			return
		}
		if !ok {
			code = http.StatusNotFound
			a.error(ctx, w, r, code, ErrNotFound)
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
