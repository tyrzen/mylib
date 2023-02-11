package rest

import (
	"context"
	"net/http"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/app/models"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Reader struct {
	logic ReaderLogic
	resp  responder
}

func NewReader(logic ReaderLogic, logger models.Logger) Reader {
	return Reader{
		logic: logic,
		resp:  responder{logger},
	}
}

func (r Reader) Route(router chi.Router) {
	router.Route("/readers", func(router chi.Router) {
		router.Post("/", r.Create)
	})

	router.Route("/auth", func(router chi.Router) {
		router.Post("/login", r.Login)
		router.With(r.resp.WithAuth).Post("/token", r.Refresh)
		router.With(r.resp.WithAuth).Post("/logout", r.Logout)
	})

}

// Create creates new models.Reader.
func (r Reader) Create(rw http.ResponseWriter, req *http.Request) {
	var reader models.Reader
	if err := r.resp.DecodeBody(req, &reader); err != nil {
		r.resp.Write(rw, req, http.StatusBadRequest, ErrDecoding)
		r.resp.Errorw("Failed decoding reader data from request.", "error", err)

		return
	}

	if err := reader.OK(); err != nil {
		r.resp.Write(rw, req, http.StatusBadRequest, err)
		r.resp.Debugw("Failed validating reader.", "error", err)

		return
	}

	reader.Normalize()

	if err := reader.HashPassword(); err != nil {
		r.resp.Write(rw, req, http.StatusInternalServerError, exceptions.ErrHashing)
		r.resp.Errorw("Failed hashing readers password.", "error", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	if err := r.logic.SignUp(ctx, reader); err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			r.resp.Write(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrDuplicateEmail):
			r.resp.Write(rw, req, http.StatusConflict, exceptions.ErrDuplicateEmail)
		case errors.Is(err, exceptions.ErrDuplicateID):
			r.resp.Write(rw, req, http.StatusConflict, exceptions.ErrDuplicateID)
		default:
			r.resp.Write(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		r.resp.Errorw("Failed creating reader.", "error", err)

		return
	}

	msg := response{Message: "Reader successfully created."}
	r.resp.Write(rw, req, http.StatusCreated, msg)
	r.resp.Debugf(msg.Message)
}

// Login handles authorization process of created models.Reader.
func (r Reader) Login(rw http.ResponseWriter, req *http.Request) {
	var creds models.Credentials
	if err := r.resp.DecodeBody(req, &creds); err != nil {
		r.resp.Write(rw, req, http.StatusBadRequest, ErrDecoding)
		r.resp.Errorw("Failed decoding reader data from request.", "error", err)

		return
	}

	creds.Normalize()

	if err := creds.OK(); err != nil {
		r.resp.Write(rw, req, http.StatusBadRequest, err)
		r.resp.Debugf("Failed validating %T: %v", creds, err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	tokenPair, err := r.logic.SignIn(ctx, creds)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			r.resp.Write(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrTokenNotCreated):
			r.resp.Write(rw, req, http.StatusBadGateway, exceptions.ErrTokenNotCreated)
		case errors.Is(err, exceptions.ErrInvalidCredits):
			r.resp.Write(rw, req, http.StatusUnauthorized, exceptions.ErrNotAuthorized)
		default:
			r.resp.Write(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		r.resp.Debugw("Failed signup reader.", "error", err)

		return
	}

	setCookie(rw, refreshTokenKey, tokenPair.RefreshToken, tokenPair.ExpiresIn, "auth")

	r.resp.Write(rw, req.WithContext(ctx), http.StatusOK, tokenPair)
	r.resp.Debugf("Reader authorized successfully.")
}

// Logout handles logout process.
func (r Reader) Logout(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	accessToken := retrieveToken[models.AccessToken](req)

	if err := r.logic.SignOut(ctx, *accessToken); err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			r.resp.Write(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrTokenExpired),
			errors.Is(err, exceptions.ErrTokenInvalid),
			errors.Is(err, exceptions.ErrTokenInvalidSigningMethod):
			r.resp.Write(rw, req, http.StatusUnauthorized, err)
		default:
			r.resp.Write(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		r.resp.Debugw("Failed signup reader.",
			accessTokenKey, accessToken,
			"error", err)

		return
	}

	setCookie(rw, refreshTokenKey, "", -1, "")

	msg := response{Message: "Reader logout successfully."}
	r.resp.Write(rw, req, http.StatusOK, msg)
	r.resp.Debugf(msg.Message)
}

// Refresh handles process of token pair refreshment.
func (r Reader) Refresh(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	refreshToken := retrieveToken[models.RefreshToken](req)

	tokenPair, err := r.logic.Refresh(ctx, *refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			r.resp.Write(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrTokenNotFound):
			r.resp.Write(rw, req, http.StatusBadRequest, exceptions.ErrTokenNotFound)
		case errors.Is(err, exceptions.ErrTokenNotCreated):
			r.resp.Write(rw, req, http.StatusBadGateway, exceptions.ErrTokenNotCreated)
		case errors.Is(err, exceptions.ErrTokenInvalid):
			r.resp.Write(rw, req, http.StatusForbidden, exceptions.ErrTokenInvalid)
		default:
			r.resp.Write(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		r.resp.Debugw("Failed refresh readers tokens.",
			refreshTokenKey, refreshToken,
			"error", err)

		return
	}

	setCookie(rw, refreshTokenKey, tokenPair.RefreshToken, tokenPair.ExpiresIn, "auth")

	r.resp.Write(rw, req.WithContext(ctx), http.StatusOK, tokenPair)
	r.resp.Debugf("Readers tokens refreshed successfully.")
}
