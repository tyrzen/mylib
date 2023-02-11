package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

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

	b.resp.Write(rw, req, http.StatusOK, book)
	b.resp.Debugf("Book fetched successfully.")
}

// FindMany handles bulk fetching books by given OData query.
// If $top number not specified maxOnPage will be used.
// In case number of requested books is more than maxOnPage
// then nextLink will be rendered in response.
func (b Book) FindMany(rw http.ResponseWriter, req *http.Request) {
	filter, err := models.NewDataFilter[models.Book](req.URL)
	if err != nil {
		b.resp.Write(rw, req, http.StatusBadRequest, ErrInvalidQuery)
		b.resp.Errorw("Failed parsing query from request URL.", "error", err)

		return
	}

	maxOnPage, err := strconv.Atoi(os.Getenv("BOOKS_MAX_ON_PAGE"))
	if err != nil {
	}

	delta := filter.Top - maxOnPage
	filter.Top = maxOnPage

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	books, err := b.logic.FetchMany(ctx, *filter)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			b.resp.Write(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
			// not sure about that.
		case errors.Is(err, exceptions.ErrRecordNotFound):
			b.resp.Write(rw, req, http.StatusBadRequest, exceptions.ErrRecordNotFound)
		default:
			b.resp.Write(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		b.resp.Errorw("Failed fetching books.", "error", err)

		return
	}

	var nextLink string
	if delta > 0 && len(books) == maxOnPage {
		filter.Skip = maxOnPage
		filter.Top = delta
		filter.UpdateURL()
		nextLink = filter.URL.String()
	}

	resp := struct {
		Books    []models.Book `json:"books"`
		NextLink string        `json:"next_link,omitempty"`
	}{
		Books:    books,
		NextLink: nextLink,
	}

	b.resp.Write(rw, req, http.StatusOK, resp)
	b.resp.Debugf("Books fetched successfully.")
}
