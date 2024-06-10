package clickhouse

import (
	"context"
	"log/slog"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	logger "github.com/rautaruukkipalich/urlsh/pkg/log"
)

func createDatabase(ctx context.Context, conn driver.Conn) error {
	const op = "store.clickhouse.CreateDatabase"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))
	// conn.Exec(ctx, "DROP TABLE IF EXIST example")

	err := conn.Exec(ctx, "CREATE DATABASE IF NOT EXISTS default")
	if err != nil {
		logger.LoggerFromContext(ctx).Info("database exist", slog.String("err", err.Error()))
	}

	err = conn.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS default.urls (
				short String,
				long String
			) engine=Memory;
		`)

	if err != nil {
		logger.LoggerFromContext(ctx).Info("table exist", slog.String("err", err.Error()))
	}
	return err
}
