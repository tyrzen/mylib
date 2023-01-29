package rest

import (
	"context"
	"net/http"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/app/models"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Book struct {
	logic BookLogic
	resp  responder
}

func NewBook(logic BookLogic, logger models.Logger) Book {
	return Book{
		logic: logic,
		resp:  responder{logger},
	}
}

func (b Book) Route(router chi.Router) {
	router.With(b.resp.WithAuth, b.resp.WithRole("admin")).
		Route("/books", func(router chi.Router) {
			router.Post("/", b.Create)
		})
}

func (b Book) Create(rw http.ResponseWriter, req *http.Request) {
	var book models.Book
	if err := b.resp.DecodeBody(req, &book); err != nil {
		b.resp.Write(rw, req, http.StatusBadRequest, ErrDecoding)
		b.resp.Errorw("Failed decoding book data from request.", "error", err)

		return
	}

	if err := book.OK(); err != nil {
		b.resp.Write(rw, req, http.StatusBadRequest, err)
		b.resp.Debugw("Failed validating book.", "error", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	if err := b.logic.Import(ctx, book); err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			b.resp.Write(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
			// TODO: finish error handling.
		default:
			b.resp.Write(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		b.resp.Errorw("Failed importing book.", "error", err)

		return
	}

	msg := response{Message: "Book imported successfully."}
	b.resp.Write(rw, req, http.StatusCreated, msg)
	b.resp.Debugf(msg.Message)
}
