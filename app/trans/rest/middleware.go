package rest

import (
	"net/http"

	"github.com/delveper/mylib/app/ent"
)

func ChainMiddlewares(h http.Handler, mds ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range mds {
		h = m(h)
	}

	return h
}

func WithLogRequest(logger ent.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			logger.Debugw("Request:",
				"Method", req.Method,
				"URL", req.URL,
				"User-Agent", req.UserAgent(),
			)

			h.ServeHTTP(rw, req)
		})
	}
}
