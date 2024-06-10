package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	logger "github.com/rautaruukkipalich/urlsh/pkg/log"
	"github.com/rautaruukkipalich/urlsh/internal/model"
)

// prepare?
func (db *DB) SetURLs(ctx context.Context, urls *model.URLs) error {
	const op = "store.clickhouse.SetURLs"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))

	defer func() {
		if r := recover(); r != nil {
			logger.LoggerFromContext(ctx).Warn("recovered", slog.Any("recovered_from", r))
		}
	}()

	if err := db.Conn.Exec(
		ctx,
		`INSERT INTO default.urls VALUES ($1, $2);`,
		urls.Short,
		urls.Long,
	); err != nil {
		logger.LoggerFromContext(ctx).Error("error while push urls", slog.Any("error", err.Error()))
		return err
	}

	return nil
}

func (db *DB) GetShortURL(ctx context.Context, urls *model.URLs) (bool, error) {
	const op = "store.clickhouse.GetShortURL"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))

	defer func() {
		if r := recover(); r != nil {
			logger.LoggerFromContext(ctx).Warn("recovered", slog.Any("recovered_from", r))
		}
	}()

	row := db.Conn.QueryRow(
		ctx,
		`SELECT short FROM default.urls WHERE long = $1;`,
		urls.Long,
	)

	if row.Err() != nil {
		logger.LoggerFromContext(ctx).Error("error while get short url", slog.Any("error", row.Err().Error()))
		return false, row.Err()
	}

	if err := row.Scan(&urls.Short); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		logger.LoggerFromContext(ctx).Error("error while scan row", slog.Any("error", err.Error()))
		return false, err
	}

	return true, nil
}

func (db *DB) GetLongURL(ctx context.Context, urls *model.URLs) (bool, error) {
	const op = "store.clickhouse.GetLongURL"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))

	defer func() {
		if r := recover(); r != nil {
			logger.LoggerFromContext(ctx).Warn("recovered", slog.Any("recovered_from", r))
		}
	}()

	row := db.Conn.QueryRow(
		ctx,
		`SELECT long FROM default.urls WHERE short = $1;`,
		urls.Short,
	)
	if row.Err() != nil {
		logger.LoggerFromContext(ctx).Error("err while get long url", slog.Any("error", row.Err().Error()))
		return false, row.Err()
	}

	if err := row.Scan(&urls.Long); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		logger.LoggerFromContext(ctx).Error("err while scan row", slog.Any("error", err.Error()))
		return false, err
	}

	return true, nil
}
