package rest

import (
	"context"
	"net/http"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/delveper/mylib/lib/tokay"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Reader struct {
	ReaderLogic
	SessionLogic
}

type ReaderKey int

const readerKey ReaderKey = iota

func NewReader(logic ReaderLogic, session SessionLogic) Reader {
	return Reader{
		ReaderLogic:  logic,
		SessionLogic: session,
	}
}

func (r Reader) Route(router chi.Router) {
	router.Route("/readers", func(router chi.Router) {
		router.Method(http.MethodPost, "/", r.Create())
	})

	router.Route("/auth", func(router chi.Router) {
		router.Method(http.MethodPost, "/login", r.Login())
		router.Method(http.MethodPost, "/logout", r.Logout())
		router.Method(http.MethodPost, "/token", nil)
	})
}

func withReader(ctx context.Context, reader *ent.Reader) context.Context {
	return context.WithValue(ctx, readerKey, reader)
}

func extractReader(ctx context.Context) *ent.Reader {
	val := ctx.Value(readerKey)
	reader, ok := val.(*ent.Reader)
	if !ok {
		return nil
	}

	return reader
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

		logger.Debugf("Reader validated.")

		reader.Normalize()
		logger.Debugf("Reader normalized.")

		if err := reader.HashPassword(); err != nil {
			respond(rw, req, http.StatusInternalServerError, exc.ErrHashing)
			logger.Errorw("Failed hashing readers password.", "error", err)

			return
		}

		logger.Debugw("Readers password hashed.")

		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		if err := r.SignUp(ctx, reader); err != nil {
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
		logger.Debugw(resp.Message)
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

		reader, err := r.SignIn(ctx, creds)
		if err != nil {
			switch {
			case errors.Is(err, exc.ErrDeadline):
				respond(rw, req, http.StatusGatewayTimeout, exc.ErrDeadline)
			case errors.Is(err, exc.ErrInvalidCredits):
				respond(rw, req, http.StatusUnauthorized, exc.ErrNotAuthorized)
			default:
				respond(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
			}

			logger.Debugw("Failed signup reader", "error", err)

			return
		}

		tokenPair, err := tokay.NewTokenPair(reader.ID, reader.IsAdmin)
		if err != nil {
			respond(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
			logger.Errorw("Failed creating token.", "error", err)

			return
		}

		token := ent.NewToken(tokenPair.Refresh.ID, reader.ID, tokenPair.Refresh.Expiration)

		if err := r.SessionLogic.Create(ctx, token); err != nil {
			respond(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
			logger.Errorw("Failed creating session.", "error", err)

			return
		}

		setCookie(rw, "refresh_token", tokenPair.Refresh.Value, tokenPair.Refresh.Expiration)

		respond(rw, req, http.StatusOK, tokenPair)
		logger.Debugf("Reader authorized successfully.")
	}
}

func (r Reader) Logout() HandlerLoggerFunc {
	return func(rw http.ResponseWriter, req *http.Request, logger ent.Logger) {
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	}
}

func (r Reader) Refresh() HandlerLoggerFunc {
	return func(rw http.ResponseWriter, req *http.Request, logger ent.Logger) {
		cookie, err := req.Cookie(tokenKey)
		if errors.Is(err, http.ErrNoCookie) {
			respond(rw, req, http.StatusUnauthorized, exc.ErrNotAuthorized)
			logger.Errorw("Failed retrieving token from cookie.", "error", err)

			return
		}

		if err != nil {
			respond(rw, req, http.StatusBadRequest, err)
			logger.Errorw("Failed retrieving token from cookie.", "error", err)

			return
		}

		if err := tokay.Check(cookie.Value); err != nil {
			switch {
			case errors.Is(err, exc.ErrTokenExpired),
				errors.Is(err, exc.ErrTokenInvalid),
				errors.Is(err, exc.ErrTokenInvalidSigningMethod):
				respond(rw, req, http.StatusUnauthorized, err)
				logger.Errorw("Failed validate token.", "error", err)
			default:
				respond(rw, req, http.StatusBadRequest, err)
				logger.Errorw("Failed VALIDATE token.", "error", exc.ErrUnexpected)
			}

			return
		}

		token, err := tokay.NewTokenPair("", true)
		if err != nil {
			respond(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
			logger.Errorw("Failed creating token.", "error", err)

			return
		}
		_ = token

		resp := response{Message: "Token refreshed successfully."}
		respond(rw, req, http.StatusOK, resp)
		logger.Debugf(resp.Message)
	}
}
