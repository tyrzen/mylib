package rest

import (
	"context"
	"net/http"
	"os"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/app/models"
	"github.com/delveper/mylib/lib/tokay"
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

// WithLogRequest logs every request and sends logger instance to further handler.
func WithLogRequest(logger models.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			var id string
			if id = req.Header.Get(xRequestID); id == "" {
				id = uuid.New().String()
			}

			req.Header.Set(xRequestID, id)
			next = WithContextKey(requestContextKey, id)(next)

			logger.Debugw("Request:",
				"id", id,
				"method", req.Method,
				"uri", req.RequestURI,
				"user-agent", req.UserAgent(),
				"remote", req.RemoteAddr,
			)

			next.ServeHTTP(rw, req)
		})
	}
}

// WithAuth will check if token is valid.
func (r responder) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		val := retrieveJWT(req)
		key := os.Getenv("JWT_KEY")

		// TODO: Make logic token-stateful.
		token, err := tokay.Parse[models.AccessToken](val, key)
		if err != nil {
			switch {
			case errors.Is(err, exceptions.ErrTokenExpired),
				errors.Is(err, exceptions.ErrTokenInvalid),
				errors.Is(err, exceptions.ErrTokenNotFound),
				errors.Is(err, exceptions.ErrTokenInvalidSigningMethod):
				r.writeJSON(rw, req, http.StatusUnauthorized, err)
			default:
				r.writeJSON(rw, req, http.StatusBadRequest, err)
			}

			r.Errorw("Failed to validate token.",
				"access_token", val,
				"error", err)

			return
		}

		r.Debugw("Token validated.", "token", token)

		next = WithContextKey(tokenContextKey, token)(next)
		next.ServeHTTP(rw, req)
	})
}

func (r responder) WithAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		token := retrieveToken[models.AccessToken](req)

		if token == nil {
			r.writeJSON(rw, req, http.StatusUnprocessableEntity, exceptions.ErrTokenNotFound)
			r.Errorw("Failed retrieve token from context.", "error", exceptions.ErrTokenNotFound)

			return
		}

		if token.Role != "admin" {
			r.writeJSON(rw, req, http.StatusUnauthorized, ErrPermissions)
			r.Infow("Failed check permissions.", "error", ErrPermissions)

			return
		}

		next.ServeHTTP(rw, req)
	})
}

// WithoutPanic recovers from panic.
func WithoutPanic(logger models.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if err := recover(); err != nil {
				logger.Errorf("Recovered from panic.")
			}

			next.ServeHTTP(rw, req)
		})
	}
}
