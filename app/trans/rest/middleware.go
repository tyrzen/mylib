package rest

import (
	"context"
	"net/http"

	"github.com/delveper/mylib/app/ent"
	"github.com/google/uuid"
)

func ChainMiddlewares(hdl http.Handler, mds ...func(http.Handler) http.Handler) http.Handler {
	for _, md := range mds {
		hdl = md(hdl)
	}

	return hdl
}

// WithContextKey wraps key/val pairs into request context.
func WithContextKey(key any, val any) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx := context.WithValue(req.Context(), key, val)
			next.ServeHTTP(rw, req.WithContext(ctx))
		})
	}
}

// WithRequestID generates request specific id if there is no any in headers.
func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var id string
		if id = req.Header.Get(RequestID); id == "" {
			id = uuid.New().String()
		}
		req.Header.Set(RequestID, id)

		ctx := context.WithValue(req.Context(), RequestID, id)

		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}

// WithLogger logs every request and sends logger instance to further handler.
func WithLogger(logger ent.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			logger.Debugw("Request:",
				"id", req.Header.Get(RequestID),
				"method", req.Method,
				"uri", req.RequestURI,
				"user-agent", req.UserAgent(),
				"remote", req.RemoteAddr,
			)

			WithContextKey(loggerKey, logger)(next).ServeHTTP(rw, req)
		})
	}
}
