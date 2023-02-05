package rest

import (
	"context"
	"fmt"
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
			router.Get("/{id}", b.Find)
			router.Get("/", b.FindMany)
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
		case errors.Is(err, exceptions.ErrRecordNotFound):
			b.resp.Write(rw, req, http.StatusExpectationFailed, fmt.Errorf("author : %w", exceptions.ErrRecordNotFound))
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

func (b Book) Find(rw http.ResponseWriter, req *http.Request) {
	var book models.Book
	book.ID = chi.URLParam(req, "id")

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	book, err := b.logic.Fetch(ctx, book)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			b.resp.Write(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrRecordNotFound):
			b.resp.Write(rw, req, http.StatusExpectationFailed, exceptions.ErrRecordNotFound)
		default:
			b.resp.Write(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		b.resp.Errorw("Failed fetching book.", "error", err)

		return
	}

	msg := response{Message: "Book fetched successfully."}
	b.resp.Write(rw, req, http.StatusOK, msg)
	b.resp.Debugf(msg.Message)
}

func (b Book) FindMany(rw http.ResponseWriter, req *http.Request) {
	var filter models.DataFilter
	if err := filter.ParseURL(req.URL.String(), models.Book{}); err != nil {
		b.resp.Write(rw, req, http.StatusBadRequest, ErrInvalidQuery)
		b.resp.Errorw("Failed parsing query from request URL.", "error", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	books, err := b.logic.FetchMany(ctx, filter)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			b.resp.Write(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
			// TODO: Add cases.
		default:
			b.resp.Write(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		b.resp.Errorw("Failed fetching books.", "error", err)

		return
	}

	// TODO: Make pagination
	_ = books
}
