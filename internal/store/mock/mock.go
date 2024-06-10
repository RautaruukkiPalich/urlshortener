package mock

import (
	"context"

	"github.com/rautaruukkipalich/urlsh/config"
	"github.com/rautaruukkipalich/urlsh/internal/model"
)

type DB struct {
	Conn map[string]string
}

func New(ctx context.Context, cfg *config.DBConfig) (*DB, error) {
	return &DB{
		Conn: make(map[string]string),
	}, nil
}

func (db *DB) SetURLs(ctx context.Context, urls *model.URLs) error {
	db.Conn[urls.Long] = urls.Short
	return nil
}

func (db *DB) GetShortURL(ctx context.Context, urls *model.URLs) (bool, error) {
	short, ok := db.Conn[urls.Long]
	urls.Short = short
	return ok, nil
}

func (db *DB) GetLongURL(ctx context.Context, urls *model.URLs) (bool, error) {
	for long, short := range db.Conn{
		if short == urls.Short  {
			urls.Long = long
			return true, nil
		} 
	}
	return false, nil
}

func (db *DB) Stop(ctx context.Context){}