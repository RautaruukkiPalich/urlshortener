package testserver

import (
	"net/http"

	"github.com/rautaruukkipalich/urlsh/internal/model"
)

var (
	URLsTestCases = []struct {
		name         string
		urls         model.URLs
		expectedCode int
	}{
		{"test1",
			model.URLs{
				Long:  "http://localhost.ru",
			},
			http.StatusOK,
		},
		{
			"test2",
			model.URLs{
				Long:  "http://localhost.ru/swagger/index.html#/urls/post_shorten",
			},
			http.StatusOK,
		}, {
			"test3",
			model.URLs{
				Long:  "https://github.com/RautaruukkiPalich?tab=repositories",
			},
			http.StatusOK,
		},
	}
)
