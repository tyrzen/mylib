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

// WithLogRequest logs every request and sends logger instance to further handler.
func WithLogRequest(logger models.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			logger.Debugw("Request:",
				"id", req.Header.Get(xRequestID),
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
func (r Reader) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		val := retrieveJWT(req)
		key := os.Getenv("JWT_KEY")

		token, err := tokay.Parse[models.AccessToken](val, key)
		if err != nil {
			switch {
			case errors.Is(err, exceptions.ErrTokenExpired),
				errors.Is(err, exceptions.ErrTokenInvalid),
				errors.Is(err, exceptions.ErrTokenNotFound),
				errors.Is(err, exceptions.ErrTokenInvalidSigningMethod):
				r.resp.Write(rw, req, http.StatusUnauthorized, err)
				r.resp.Errorw("Failed validate token.",
					"access_token", val,
					"error", err)
			default:
				r.resp.Write(rw, req, http.StatusBadRequest, err)
				r.resp.Errorw("Failed validate token.",
					"access_token", val,
					"error", exceptions.ErrUnexpected)
			}

			r.resp.Debugw("Token validated.", "token", token)
			return
		}

		next = WithContextKey(tokenContextKey, token)(next)
		next.ServeHTTP(rw, req)
	})
}

func (r Reader) WithRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			val := req.Context().Value(tokenContextKey)

			token, ok := val.(models.AccessToken)
			if !ok {
				r.resp.Write(rw, req, http.StatusUnprocessableEntity, exceptions.ErrTokenNotFound)
				r.resp.Errorw("Failed retrieve token from context.", "error", exceptions.ErrTokenNotFound)
			}

			if token.Role != role {
				r.resp.Write(rw, req, http.StatusUnauthorized, ErrPermissions)
				r.resp.Infow("Failed check permissions.", "error", ErrPermissions)
			}

			next.ServeHTTP(rw, req)
		})
	}
}
