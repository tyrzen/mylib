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
	ReaderLogic
}

func NewReader(logic ReaderLogic) Reader {
	return Reader{ReaderLogic: logic}
}

func (r Reader) Route(router chi.Router) {
	router.Method(http.MethodPost, "/readers", r.Create())
}

func (r Reader) Create() HandlerLoggerFunc {
	return func(rw http.ResponseWriter, req *http.Request, logger ent.Logger) {
		var reader ent.Reader
		if err := decodeBody(req, &reader); err != nil {
			respond(rw, req, http.StatusBadRequest, ErrDecoding)
			logger.Errorf("Failed decoding reader data from request.", "request", req, "error", err)

			return
		}

		if err := reader.OK(); err != nil {
			respond(rw, req, http.StatusBadRequest, err)
			logger.Debugf("Failed validating %T: %v", reader, err)

			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		err := r.SignUp(ctx, reader)
		if err != nil {
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
