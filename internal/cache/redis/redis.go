package rediscache

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/go-redis/redis"
	"github.com/rautaruukkipalich/urlsh/config"
	"github.com/rautaruukkipalich/urlsh/internal/cache"
	"github.com/rautaruukkipalich/urlsh/internal/model"
	logger "github.com/rautaruukkipalich/urlsh/pkg/log"
)

type Cache struct{
	store *redis.Client
	exp time.Duration
}

func New(ctx context.Context, cfg *config.CacheConfig) (*Cache, error) {
	const op = "internal.cache.redis.New"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))
	log := logger.LoggerFromContext(ctx)

	opts, err := redis.ParseURL(cfg.URI)
    if err != nil {
		log.Error("error", slog.Any("error", err.Error()))
        return nil, err
    }
	cache := redis.NewClient(opts)
	if err := cache.Ping().Err(); err != nil {
		log.Error("error", slog.Any("error", err.Error()))
		return nil, err
	}
	return &Cache{
		store: cache,
		exp: cfg.Exp,
	}, nil
}

func (c *Cache) GetLongURL(ctx context.Context, urls *model.URLs) (bool, error) {
	const op = "internal.cache.redis.GetLongURL"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))
	log := logger.LoggerFromContext(ctx)

	data, err := c.store.Get(urls.Short).Result()
	if err != nil {
		if err == redis.Nil {
			return false, cache.ErrNotFound
		}
		log.Error("error get long url from cache", slog.Any("error", err.Error()))
		return false, err
	}

	if err := json.Unmarshal([]byte(data), &urls.Long); err != nil {
		log.Error("error unmarshall from cache", slog.Any("error", err.Error()))
		return false, err
	}

	return true, nil
}

func (c *Cache) GetShortURL(ctx context.Context, urls *model.URLs) (bool, error) {
	const op = "internal.cache.redis.GetShortURL"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))
	log := logger.LoggerFromContext(ctx)

	data, err := c.store.Get(urls.Long).Result()
	if err != nil {
		if err == redis.Nil {
			return false, cache.ErrNotFound
		}
		log.Error("error get short url from cache", slog.Any("error", err.Error()))
		return false, err
	}

	if err := json.Unmarshal([]byte(data), &urls.Short); err != nil {
		log.Error("error unmarshall from cache", slog.Any("error", err.Error()))
		return false, err
	}

	return true, nil
}


func (c *Cache) SetURLs(ctx context.Context, urls *model.URLs) error {
	const op = "internal.cache.redis.SetURLs"
	ctx = logger.AddTrace(ctx, slog.Any("op", op))
	log := logger.LoggerFromContext(ctx)

	json_data_short, err := json.Marshal(&urls.Short)
	if err != nil {
		log.Error("error marshall to cache", slog.Any("error", err.Error()))
		return err
	}
	json_data_long, err := json.Marshal(&urls.Long)
	if err != nil {
		log.Error("error marshall to cache", slog.Any("error", err.Error()))
		return err
	}

	if err := c.store.Set(urls.Short, json_data_long, c.exp).Err(); err != nil {
		log.Error("error set short url to cache", slog.Any("error", err.Error()))
		return err
	}

	if err := c.store.Set(urls.Long, json_data_short, c.exp).Err(); err != nil {
		log.Error("error set long url to cache", slog.Any("error", err.Error()))
		return err
	}

	return nil 
}

func (c *Cache) Stop(context.Context) {
	c.store.Close()
}