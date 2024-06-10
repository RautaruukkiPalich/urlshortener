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
				Short: "123qwe11",
			},
			http.StatusOK,
		},
		{
			"test2",
			model.URLs{
				Long:  "http://localhost.ru/swagger/index.html#/urls/post_shorten",
				Short: "1234325g",
			},
			http.StatusOK,
		}, {
			"test3",
			model.URLs{
				Long:  "https://github.com/RautaruukkiPalich?tab=repositories",
				Short: "bgfj!fg",
			},
			http.StatusOK,
		},
	}
)
