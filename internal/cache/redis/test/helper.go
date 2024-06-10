package redistest

import (
	"time"

	"github.com/rautaruukkipalich/urlsh/config"
	"github.com/rautaruukkipalich/urlsh/internal/model"
)

var (
	cfg = config.CacheConfig{
		URI: "redis://localhost:6381",
		Exp: 30 * time.Second,
	}

	testURLs = []model.URLs{
		{
			Long: "http://www.clickhouse.com",
			Short: "321qwe",
		},{
			Long: "http://gmail.com",
			Short: "123321",
		},
	}
)