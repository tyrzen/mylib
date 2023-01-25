package rest

import (
	"context"
	"net/http"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func ChainMiddlewares(hdl http.Handler, mds ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mds) - 1; i >= 0; i-- { // LIFO order
		hdl = mds[i](hdl)
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

		if id = req.Header.Get(xRequestID); id == "" {
			id = uuid.New().String()
		}

		req.Header.Set(xRequestID, id)

		next = WithContextKey(requestContextKey, id)(next)
		next.ServeHTTP(rw, req)
	})
}

// WithLogger logs every request and sends logger instance to further handler.
func WithLogger(logger ent.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			logger.Debugw("Request:",
				"id", req.Header.Get(xRequestID),
				"method", req.Method,
				"uri", req.RequestURI,
				"user-agent", req.UserAgent(),
				"remote", req.RemoteAddr,
			)

			next = WithContextKey(loggerContextKey, logger)(next)
			next.ServeHTTP(rw, req)
		})
	}
}

// WithAuth will check if token is valid.
func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		logger := loggerFromContext(rw, req)

		token := tokenFromRequest(req)

		if _, err := tokay.Verify(token); err != nil {
			switch {
			case errors.Is(err, exc.ErrTokenExpired),
				errors.Is(err, exc.ErrTokenInvalid),
				errors.Is(err, exc.ErrTokenNotFound),
				errors.Is(err, exc.ErrTokenInvalidSigningMethod):
				respond(rw, req, http.StatusUnauthorized, err)
				logger.Errorw("Failed validate token.", "error", err)
			default:
				respond(rw, req, http.StatusBadRequest, err)
				logger.Errorw("Failed validate token.", "error", exc.ErrUnexpected)
			}

			return
		}

		next = WithContextKey(tokenContextKey, token)(next)
		next.ServeHTTP(rw, req)
	})
}
