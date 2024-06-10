package mockcachetest

import (
	"github.com/rautaruukkipalich/urlsh/config"
	"github.com/rautaruukkipalich/urlsh/internal/model"
)

var (
	cfg = config.CacheConfig{}

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