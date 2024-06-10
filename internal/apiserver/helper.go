package apiserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"regexp"

	logger "github.com/rautaruukkipalich/urlsh/pkg/log"
)

const (
	Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@%^*+-_"
	URLlength = 8
)

func (a *APIServer) error(ctx context.Context, w http.ResponseWriter, r *http.Request, err errorResponse) {
	a.JSONrespond(ctx, w, r, err.Code, err)
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


// func (a *APIServer) redirect(w http.ResponseWriter, r *http.Request, url string) {
// 	logger.LoggerFromContext(a.ctx).With(
// 		slog.String("remote addr", r.RemoteAddr),
// 		slog.String("from", r.RequestURI),
// 		slog.String("to", url),
// 	).Info("reditected")
// 	// a.log.With(
// 	// 	slog.String("remote addr", r.RemoteAddr),
// 	// 	slog.String("from", r.RequestURI),
// 	// 	slog.String("to", url),
// 	// ).Info("reditected")
// 	http.Redirect(w, r, url, http.StatusPermanentRedirect)
// }
