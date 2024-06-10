package clickhouse

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/rautaruukkipalich/urlsh/config"
	logger "github.com/rautaruukkipalich/urlsh/pkg/log"
)

type DB struct {
	Conn driver.Conn
}

func New(ctx context.Context, cfg *config.DBConfig) (*DB, error) {
	const op = "internal.store.clickhouse.New"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))
	log := logger.LoggerFromContext(ctx)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.URI}, //"127.0.0.1:8123"
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": cfg.MaxExecutionTime,
		},
		DialTimeout: cfg.DialTimeout,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Protocol: clickhouse.HTTP,
	})

	if err != nil {
		log.Error("error connecting to clickhouse", slog.Any("error", err))
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Error(
				"exception", 
				slog.Any("error", fmt.Sprintf("exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)),
			)
		}
		log.Error("error connecting to clickhouse", slog.Any("error", err.Error()))
		return nil, err
	}

	if err := createDatabase(ctx, conn); err != nil {
		//nolint:all
		err = nil 
	}

	return &DB{Conn: conn}, nil
}

func (db *DB) Stop(ctx context.Context) {
	const op = "store.clickhouse.Stop"
	ctx = logger.AddTrace(ctx, slog.String("op", op))
	logger.LoggerFromContext(ctx).Info("stop clickhouse")
	defer db.Conn.Close()
}
