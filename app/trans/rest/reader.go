package rest

import (
	"context"
	"net/http"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Reader struct {
	logic ReaderLogic
}

func NewReader(logic ReaderLogic) Reader {
	return Reader{
		logic: logic,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func newTokenPair(pair *ent.TokenPair) TokenPair {
	return TokenPair{
		AccessToken:  pair.Access.Value,
		RefreshToken: pair.Refresh.Value,
	}
}

func (r Reader) Route(router chi.Router) {
	router.Route("/readers", func(router chi.Router) {
		router.Method(http.MethodPost, "/", r.Create())
	})

	router.Route("/auth", func(router chi.Router) {
		router.Method(http.MethodPost, "/login", r.Login())

		router.Use(WithAuth)
		router.Method(http.MethodPost, "/logout", r.Logout())
		router.Method(http.MethodPost, "/refresh-token", nil)
	})
}

// Create creates new ent.Reader.
func (r Reader) Create() HandlerLoggerFunc {
	return func(rw http.ResponseWriter, req *http.Request, logger ent.Logger) {
		var reader ent.Reader
		if err := decodeBody(req, &reader); err != nil {
			respond(rw, req, http.StatusBadRequest, ErrDecoding)
			logger.Errorw("Failed decoding reader data from request.", "error", err)

			return
		}

		if err := reader.OK(); err != nil {
			respond(rw, req, http.StatusBadRequest, err)
			logger.Debugf("Failed validating %T: %v", reader, err)

			return
		}

		reader.Normalize()

		if err := reader.HashPassword(); err != nil {
			respond(rw, req, http.StatusInternalServerError, exc.ErrHashing)
			logger.Errorw("Failed hashing readers password.", "error", err)

			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		if err := r.logic.SignUp(ctx, reader); err != nil {
			switch {
			case errors.Is(err, exc.ErrDeadline):
				respond(rw, req, http.StatusGatewayTimeout, exc.ErrDeadline)
			case errors.Is(err, exc.ErrDuplicateEmail):
				respond(rw, req, http.StatusConflict, exc.ErrDuplicateEmail)
			case errors.Is(err, exc.ErrDuplicateID):
				respond(rw, req, http.StatusConflict, exc.ErrDuplicateID)
			default:
				respond(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
			}

			logger.Errorw("Failed creating reader.", "error", err)

			return
		}

		resp := response{Message: "Reader successfully created."}
		respond(rw, req, http.StatusCreated, resp)
		logger.Debugf(resp.Message)
	}
}

// Login handles authorization process of created ent.Reader.
func (r Reader) Login() HandlerLoggerFunc {
	return func(rw http.ResponseWriter, req *http.Request, logger ent.Logger) {
		var creds ent.Credentials
		if err := decodeBody(req, &creds); err != nil {
			respond(rw, req, http.StatusBadRequest, ErrDecoding)
			logger.Errorw("Failed decoding reader data from request.", "error", err)

			return
		}

		creds.Normalize()

		if err := creds.OK(); err != nil {
			respond(rw, req, http.StatusBadRequest, err)
			logger.Debugf("Failed validating %T: %v", creds, err)

			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		tokenPair, err := r.logic.SignIn(ctx, creds)
		if err != nil {
			switch {
			case errors.Is(err, exc.ErrDeadline):
				respond(rw, req, http.StatusGatewayTimeout, exc.ErrDeadline)
			case errors.Is(err, exc.ErrTokenCreating):
				respond(rw, req, http.StatusBadGateway, exc.ErrTokenCreating)
			case errors.Is(err, exc.ErrInvalidCredits):
				respond(rw, req, http.StatusUnauthorized, exc.ErrNotAuthorized)
			default:
				respond(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
			}

			logger.Debugw("Failed signup reader.", "error", err)

			return
		}

		setCookie(rw, tokenCookieKey, tokenPair.Refresh.Value, tokenPair.Refresh.Expiry, 0)

		resp := newTokenPair(tokenPair)
		respond(rw, req.WithContext(ctx), http.StatusOK, resp)
		logger.Debugf("Reader authorized successfully.")
	}
}

func (r Reader) Logout() HandlerLoggerFunc {
	return func(rw http.ResponseWriter, req *http.Request, logger ent.Logger) {
		token := tokenFromRequest(req)

		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		if err := r.logic.SignOut(ctx, token); err != nil {
			switch {
			case errors.Is(err, exc.ErrDeadline):
				respond(rw, req, http.StatusGatewayTimeout, exc.ErrDeadline)
			case errors.Is(err, exc.ErrTokenExpired),
				errors.Is(err, exc.ErrTokenInvalid),
				errors.Is(err, exc.ErrTokenInvalidSigningMethod):
				respond(rw, req, http.StatusUnauthorized, err)
			default:
				respond(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
			}

			logger.Debugw("Failed signup reader.", "error", err)
			return
		}

		resp := response{Message: "Reader logout went successfully."}
		respond(rw, req, http.StatusOK, resp)
		logger.Debugf(resp.Message)
	}
}

func (r Reader) Refresh() HandlerLoggerFunc {
	return func(rw http.ResponseWriter, req *http.Request, logger ent.Logger) {
		token := tokenFromRequest(req)

		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		tokenPair, err := r.logic.Refresh(ctx, token)
		if err != nil {
			respond(rw, req, http.StatusUnauthorized, err)
			logger.Errorw("Failed checking auth details during refreshing token.", "error", err)

			return
		}

		setCookie(rw, tokenCookieKey, "", -1, 0)

		respond(rw, req, http.StatusOK, tokenPair)
		logger.Debugf("Tokens refreshed successfully.")
	}
}
