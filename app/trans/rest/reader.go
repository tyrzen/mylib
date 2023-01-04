package rest

import (
	"context"
	"net/http"

	"github.com/delveper/mylib/app/ent"
	"github.com/go-chi/chi/v5"
)

type Reader struct {
	logic  ReaderLogic
	logger ent.Logger
}

func NewReader(logic ReaderLogic, logger ent.Logger) Reader {
	return Reader{
		logic:  logic,
		logger: logger,
	}
}

func (r Reader) Route(router chi.Router) {
	router.Post("/readers", r.Create())
}

func (r Reader) Create() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var reader ent.Reader
		if err := decodeBody(req, &reader); err != nil {
			respond(rw, req, Response{Message: MsgBadRequest, Details: err.Error()}, http.StatusBadRequest, r.logger)
			r.logger.Errorf("Failed decoding reader data from request.", "request", req, "error", err)

			return
		}

		if err := reader.Validate(); err != nil {
			r.logger.Debugf("Failed validating reader: %v", err)
			respond(rw, req, Response{Message: MsgBadRequest, Details: err.Error()}, http.StatusBadRequest, r.logger)

			return
		}

		ctx := context.Background() // TODO: Make context with deadlines.

		switch err := r.logic.SignUp(ctx, reader); {
		case err != nil:
			respond(rw, req, Response{Message: MsgBadRequest, Details: err.Error()}, http.StatusBadRequest, r.logger)
		}

	}
}
