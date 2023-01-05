package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/go-chi/chi/v5"
)

type Reader struct {
	logic  ReaderLogic
	logger ent.Logger
	// resp responder
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
		resp := newResponse(rw, req, r.logger)
		defer resp.Send()

		var reader ent.Reader
		if err := decodeBody(req, &reader); err != nil {
			resp.Set(http.StatusBadRequest, Message{Message: MsgBadRequest, Details: err.Error()})
			r.logger.Errorf("Failed decoding reader data from request.", "request", req, "error", err)

			return
		}

		if err := reader.Validate(); err != nil {
			resp.Set(http.StatusBadRequest, Message{Message: MsgBadRequest, Details: err.Error()})
			r.logger.Debugf("Failed validating reader: %v", err)

			return
		}

		err := r.logic.SignUp(context.Background(), reader)
		switch {
		case err == nil:
			resp.Set(http.StatusCreated, Message{Message: MsgSuccess})
			r.logger.Debugw("Reader successfully created")

			return
		case errors.Is(err, exc.ErrDuplicateEmail):
			resp.Set(http.StatusConflict, Message{Message: MsgConflict, Details: exc.ErrDuplicateEmail.Error()})
		case errors.Is(err, exc.ErrDuplicateID):
			resp.Set(http.StatusConflict, Message{Message: MsgConflict, Details: exc.ErrDuplicateID.Error()})
		default:
			resp.Set(http.StatusInternalServerError, Message{Message: MsgInternalSeverErr})
		}

		r.logger.Errorf("Failed creating reader: %+v", err)

		return
	}
}
