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

func (r Reader) Route(rtr chi.Router) {
	rtr.Route("/readers", func(rtr chi.Router) {
		rtr.Post("/signup", r.Register)
		rtr.Post("/login", r.Login)
		rtr.With(r.resp.WithAuth).Post("/token", r.Refresh)
		rtr.With(r.resp.WithAuth).Post("/logout", r.Logout)
	})
}

// Register creates new models.Reader.
func (r Reader) Register(rw http.ResponseWriter, req *http.Request) {
	var reader models.Reader
	if err := r.resp.decodeBody(req, &reader); err != nil {
		r.resp.writeJSON(rw, req, http.StatusBadRequest, ErrDecoding)
		r.resp.Errorw("Failed decoding reader data from request.", "error", err)

		return
	}

	if err := reader.OK(); err != nil {
		r.resp.writeJSON(rw, req, http.StatusBadRequest, err)
		r.resp.Debugw("Failed validating reader.", "error", err)

		return
	}

	reader.Normalize()

	if err := reader.HashPassword(); err != nil {
		r.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrHashing)
		r.resp.Errorw("Failed hashing readers password.", "error", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	if err := r.logic.SignUp(ctx, reader); err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			r.resp.writeJSON(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrDuplicateEmail):
			r.resp.writeJSON(rw, req, http.StatusConflict, exceptions.ErrDuplicateEmail)
		case errors.Is(err, exceptions.ErrDuplicateID):
			r.resp.writeJSON(rw, req, http.StatusConflict, exceptions.ErrDuplicateID)
		default:
			r.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		r.resp.Errorw("Failed creating reader.", "error", err)

		return
	}

	msg := response{Message: "Reader successfully created."}
	r.resp.writeJSON(rw, req, http.StatusCreated, msg)
	r.resp.Debugf(msg.Message)
}

// Login handles authorization process of created models.Reader.
func (r Reader) Login(rw http.ResponseWriter, req *http.Request) {
	var creds models.Credentials
	if err := r.resp.decodeBody(req, &creds); err != nil {
		r.resp.writeJSON(rw, req, http.StatusBadRequest, ErrDecoding)
		r.resp.Errorw("Failed decoding reader data from request.", "error", err)

		return
	}

	creds.Normalize()

	if err := creds.OK(); err != nil {
		r.resp.writeJSON(rw, req, http.StatusBadRequest, err)
		r.resp.Debugf("Failed validating %T: %v", creds, err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	tokenPair, err := r.logic.SignIn(ctx, creds)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			r.resp.writeJSON(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrTokenNotCreated):
			r.resp.writeJSON(rw, req, http.StatusBadGateway, exceptions.ErrTokenNotCreated)
		case errors.Is(err, exceptions.ErrRecordNotFound),
			errors.Is(err, exceptions.ErrInvalidCredits):
			r.resp.writeJSON(rw, req, http.StatusUnauthorized, ErrNotAuthorized)
		default:
			r.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		r.resp.Debugw("Failed signup reader.", "error", err)

		return
	}

	setCookie(rw, refreshTokenKey, tokenPair.RefreshToken, tokenPair.ExpiresIn, "auth")

	r.resp.writeJSON(rw, req.WithContext(ctx), http.StatusOK, tokenPair)
	r.resp.Debugf("Reader authorized successfully.")
}

// Logout handles logout process.
func (r Reader) Logout(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	accessToken := retrieveToken[models.AccessToken](req)
	if accessToken == nil {
		r.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		r.resp.Errorf("Failed retrieve token from context.")

		return
	}

	if err := r.logic.SignOut(ctx, *accessToken); err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			r.resp.writeJSON(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrTokenExpired),
			errors.Is(err, exceptions.ErrTokenInvalid),
			errors.Is(err, exceptions.ErrTokenInvalidSigningMethod):
			r.resp.writeJSON(rw, req, http.StatusUnauthorized, err)
		default:
			r.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		r.resp.Debugw("Failed signup reader.",
			accessTokenKey, accessToken,
			"error", err)

		return
	}

	setCookie(rw, refreshTokenKey, "", -1, "")

	msg := response{Message: "Reader logout successfully."}
	r.resp.writeJSON(rw, req, http.StatusOK, msg)
	r.resp.Debugf(msg.Message)
}

// Refresh handles process of token pair refreshment.
func (r Reader) Refresh(rw http.ResponseWriter, req *http.Request) {
	var refreshToken models.RefreshToken

	if err := r.resp.decodeBody(req, &refreshToken); err != nil {
		r.resp.writeJSON(rw, req, http.StatusBadRequest, ErrDecoding)
		r.resp.Errorw("Failed decoding refresh token from request.", "error", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	tokenPair, err := r.logic.Refresh(ctx, refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			r.resp.writeJSON(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrTokenNotFound):
			r.resp.writeJSON(rw, req, http.StatusBadRequest, exceptions.ErrTokenNotFound)
		case errors.Is(err, exceptions.ErrTokenNotCreated):
			r.resp.writeJSON(rw, req, http.StatusBadGateway, exceptions.ErrTokenNotCreated)
		case errors.Is(err, exceptions.ErrTokenInvalid):
			r.resp.writeJSON(rw, req, http.StatusForbidden, exceptions.ErrTokenInvalid)
		default:
			r.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		r.resp.Debugw("Failed refresh readers tokens.",
			refreshTokenKey, refreshToken,
			"error", err)

		return
	}

	setCookie(rw, refreshTokenKey, tokenPair.RefreshToken, tokenPair.ExpiresIn, "auth")

	r.resp.writeJSON(rw, req.WithContext(ctx), http.StatusOK, tokenPair)
	r.resp.Debugf("Readers tokens refreshed successfully.")
}
