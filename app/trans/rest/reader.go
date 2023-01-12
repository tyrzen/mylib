package rest

import (
	"context"
	"net/http"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Reader struct{ ReaderLogic }

type ReaderSinger struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewReader(logic ReaderLogic) Reader {
	return Reader{ReaderLogic: logic}
}

func (r Reader) Route(router chi.Router) {
	router.Route("/readers", func(router chi.Router) {
		router.Method(http.MethodPost, "/", r.Create())
		// router.Method(http.MethodPost, "/readers/", r.Authenticate())
	})

}

// Create creates new ent.Reader.
func (r Reader) Create() HandlerLoggerFunc {
	return func(rw http.ResponseWriter, req *http.Request, logger ent.Logger) {
		var reader ent.Reader
		if err := decodeBody(req, &reader); err != nil {
			respond(rw, req, http.StatusBadRequest, ErrDecoding)
			logger.Errorw("Failed decoding reader data from request.",
				"request", req,
				"error", err)

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
			logger.Errorf("Failed hashing readers password: %+v", err)

			return
		}

		logger.Debugw("Readers password hashed.")

		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		if err := r.SignUp(ctx, reader); err != nil {
			switch {
			case errors.Is(err, exc.ErrDuplicateEmail):
				respond(rw, req, http.StatusConflict, exc.ErrDuplicateEmail)
			case errors.Is(err, exc.ErrDuplicateID):
				respond(rw, req, http.StatusConflict, exc.ErrDuplicateID)
			case errors.Is(err, exc.ErrDeadline):
				respond(rw, req, http.StatusInternalServerError, exc.ErrDeadline)
			default:
				respond(rw, req, http.StatusInternalServerError, exc.ErrUnexpected)
			}

			logger.Errorf("Failed creating reader: %+v", err)

			return
		}

		respond(rw, req, http.StatusCreated, Response{Message: "Success"})
		logger.Debugw("Reader successfully created")
	}
}

// Authenticate logins existing ent.Reader.
func (r Reader) Authenticate() HandlerLoggerFunc {
	return func(rw http.ResponseWriter, req *http.Request, logger ent.Logger) {

	}
}
