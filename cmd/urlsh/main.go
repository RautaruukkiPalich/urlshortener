package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rautaruukkipalich/urlsh/internal/apiserver"
	rediscache "github.com/rautaruukkipalich/urlsh/internal/cache/redis"
	"github.com/rautaruukkipalich/urlsh/internal/metrics"
	"github.com/rautaruukkipalich/urlsh/internal/store/clickhouse"
	"github.com/rautaruukkipalich/urlsh/pkg/configloader"
	logger "github.com/rautaruukkipalich/urlsh/pkg/log"
	"github.com/rautaruukkipalich/urlsh/pkg/loggerloader"
)

func main() {
	ctx := context.TODO()
	const op = "cmd.urlsh.main"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))
	
	cfg := configloader.MustLoadConfig()
	log := loggerloader.MustRunLogger(cfg.LogEnv)

	logger.CreateGlobalLogger(log)
	
	ctx = logger.ContextWithLogger(ctx, log)

	db, err := clickhouse.New(ctx, &cfg.Database)
	if err != nil {
		logger.LoggerFromContext(ctx).Warn("error creating database", slog.Any("error", err.Error()))
		panic(err)
	}

	cache, err := rediscache.New(ctx, &cfg.Cache)
	if err != nil {
		logger.LoggerFromContext(ctx).Warn("error creating cache store", slog.Any("error", err.Error()))
		panic(err)
	}

	srv := apiserver.New(
		ctx, 
		&cfg.Server, 
		db,
		db,
		db, 
		cache,
		cache,
		cache,
	)

	logger.LoggerFromContext(ctx).Info("starting server...")

	//run metrics
	go func(){
		_ = metrics.Listen(cfg.Metrics.URI)
	}()

	//run server
	go srv.MustRun(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	// stop server
	if err := srv.Stop(ctx); err != nil {
		logger.LoggerFromContext(ctx).Info("error stopping server: %v", err)
	}

	logger.LoggerFromContext(ctx).Warn("server stopped")

}
