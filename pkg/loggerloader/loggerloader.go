package loggerloader

import (
	"log/slog"
	"os"

	"github.com/rautaruukkipalich/prettyslog"
)

const (
	LOCAL = "local"
	DEV = "dev"
	PROD = "prod"
)

func MustRunLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case LOCAL:
		log = prettyslog.NewPrettyLogger("\t")
	case DEV:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, 
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
					ReplaceAttr: replaceAttr,
				}),
		)
	case PROD:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, 
				&slog.HandlerOptions{
					Level: slog.LevelWarn,
					ReplaceAttr: replaceAttr,
			}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, 
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
					ReplaceAttr: replaceAttr,
				}),
		)
		log.Error("invalid logger env")
		panic("invalid logger env")
	}
	return log
}