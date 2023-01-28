package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Reader struct {
	logic ReaderLogic
	resp  responder
}

func NewReader(logic ReaderLogic, logger ent.Logger) Reader {
	return Reader{
		logic: logic,
		resp:  responder{logger},
	}
}

type TokenPair struct {
	AccessToken  string        `json:"access_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    time.Duration `json:"expires_in"`
	RefreshToken string        `json:"refresh_token"`
}

func newTokenPair(pair *ent.TokenPair) TokenPair {
	return TokenPair{
		AccessToken:  pair.Access.Value,
		TokenType:    "bearer",
		ExpiresIn:    pair.Access.Expiry,
		RefreshToken: pair.Refresh.Value,
	}
}

func (r Reader) Route(router chi.Router) {
	router.Route("/auth", func(router chi.Router) {
		router.Post("/token", r.Refresh)
		router.Post("/login", r.Login)
		router.With(r.WithAuth).Post("/logout", r.Logout)
	})

	router.Route("/readers", func(router chi.Router) {
		router.Post("/", r.Create)
	})
}

// Create creates new ent.Reader.
func (r Reader) Create(rw http.ResponseWriter, req *http.Request) {
	var reader ent.Reader
	if err := r.resp.decodeBody(req, &reader); err != nil {
		r.resp.writeResponse(rw, req, http.StatusBadRequest, ErrDecoding)
		r.resp.Errorw("Failed decoding reader data from request.", "error", err)

		return
	}

	if err := reader.OK(); err != nil {
		r.resp.writeResponse(rw, req, http.StatusBadRequest, err)
		r.resp.Debugf("Failed validating %T: %v", reader, err)

		return
	}

	reader.Normalize()

	if err := reader.HashPassword(); err != nil {
		r.resp.writeResponse(rw, req, http.StatusInternalServerError, exc.ErrHashing)
		r.resp.Errorw("Failed hashing readers password.", "error", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	if err := r.logic.SignUp(ctx, reader); err != nil {
		switch {
		case errors.Is(err, exc.ErrDeadline):
			r.resp.writeResponse(rw, req, http.StatusGatewayTimeout, exc.ErrDeadline)
		case errors.Is(err, exc.ErrDuplicateEmail):
			r.resp.writeResponse(rw, req, http.StatusConflict, exc.ErrDuplicateEmail)
		case errors.Is(err, exc.ErrDuplicateID):
			r.resp.writeResponse(rw, req, http.StatusConflict, exc.ErrDuplicateID)
		default:
			r.resp.writeResponse(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
		}

		r.resp.Errorw("Failed creating reader.", "error", err)

		return
	}

	response := response{Message: "Reader successfully created."}
	r.resp.writeResponse(rw, req, http.StatusCreated, response)
	r.resp.Debugf(response.Message)
}

// Login handles authorization process of created ent.Reader.
func (r Reader) Login(rw http.ResponseWriter, req *http.Request) {
	var creds ent.Credentials
	if err := r.resp.decodeBody(req, &creds); err != nil {
		r.resp.writeResponse(rw, req, http.StatusBadRequest, ErrDecoding)
		r.resp.Errorw("Failed decoding reader data from request.", "error", err)

		return
	}

	creds.Normalize()

	if err := creds.OK(); err != nil {
		r.resp.writeResponse(rw, req, http.StatusBadRequest, err)
		r.resp.Debugf("Failed validating %T: %v", creds, err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	tokenPair, err := r.logic.SignIn(ctx, creds)
	if err != nil {
		switch {
		case errors.Is(err, exc.ErrDeadline):
			r.resp.writeResponse(rw, req, http.StatusGatewayTimeout, exc.ErrDeadline)
		case errors.Is(err, exc.ErrTokenCreating):
			r.resp.writeResponse(rw, req, http.StatusBadGateway, exc.ErrTokenCreating)
		case errors.Is(err, exc.ErrInvalidCredits):
			r.resp.writeResponse(rw, req, http.StatusUnauthorized, exc.ErrNotAuthorized)
		default:
			r.resp.writeResponse(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
		}

		r.resp.Debugw("Failed signup reader.", "error", err)

		return
	}

	setCookie(rw, refreshTokenKey, tokenPair.Refresh.Value, tokenPair.Refresh.Expiry, tokenPair.Access.Expiry)

	resp := newTokenPair(tokenPair)
	r.resp.writeResponse(rw, req.WithContext(ctx), http.StatusOK, resp)
	r.resp.Debugf("Reader authorized successfully.")
}

// Logout handles logout process.
func (r Reader) Logout(rw http.ResponseWriter, req *http.Request) {
	token := retrieveToken(req)

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	if err := r.logic.SignOut(ctx, token); err != nil {
		switch {
		case errors.Is(err, exc.ErrDeadline):
			r.resp.writeResponse(rw, req, http.StatusGatewayTimeout, exc.ErrDeadline)
		case errors.Is(err, exc.ErrTokenExpired),
			errors.Is(err, exc.ErrTokenInvalid),
			errors.Is(err, exc.ErrTokenInvalidSigningMethod):
			r.resp.writeResponse(rw, req, http.StatusUnauthorized, err)
		default:
			r.resp.writeResponse(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
		}

		r.resp.Debugw("Failed signup reader.", "error", err)

		return
	}

	setCookie(rw, refreshTokenKey, "", -1, 0)

	resp := response{Message: "Reader logout went successfully."}
	r.resp.writeResponse(rw, req, http.StatusOK, resp)
	r.resp.Debugf(resp.Message)
}

// Refresh handles process of token pair refreshment.
func (r Reader) Refresh(rw http.ResponseWriter, req *http.Request) {
	token := retrieveToken(req)

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	tokenPair, err := r.logic.Refresh(ctx, token)
	if err != nil {
		switch {
		case errors.Is(err, exc.ErrDeadline):
			r.resp.writeResponse(rw, req, http.StatusGatewayTimeout, exc.ErrDeadline)
		case errors.Is(err, exc.ErrTokenNotFound):
			r.resp.writeResponse(rw, req, http.StatusBadRequest, exc.ErrTokenNotFound)
		case errors.Is(err, exc.ErrTokenCreating):
			r.resp.writeResponse(rw, req, http.StatusBadGateway, exc.ErrTokenCreating)
		case errors.Is(err, exc.ErrTokenInvalid):
			r.resp.writeResponse(rw, req, http.StatusForbidden, exc.ErrTokenInvalid)
		default:
			r.resp.writeResponse(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
		}

		r.resp.Debugw("Failed refresh readers tokens.", "error", err)

		return
	}

	setCookie(rw, refreshTokenKey, tokenPair.Refresh.Value, tokenPair.Refresh.Expiry, tokenPair.Access.Expiry)

	resp := newTokenPair(tokenPair)
	r.resp.writeResponse(rw, req.WithContext(ctx), http.StatusOK, resp)
	r.resp.Debugf("Readers tokens refreshed successfully.")
}
