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
	router.With(b.resp.WithAuth).
		Route("/books", func(router chi.Router) {
			router.Get("/{id}", b.Find)
			router.Get("/", b.FindMany)
			router.With(b.resp.WithRole("admin")).Post("/", b.Create)
		})

	router.With(b.resp.WithAuth).
		Route("/readers/me", func(router chi.Router) {
			router.Post("/favorites", b.AddToFavorites)
			router.Post("/wishlist", b.AddToWishlist)
		})
}

func (b Book) Create(rw http.ResponseWriter, req *http.Request) {
	var book models.Book
	if err := b.resp.decodeBody(req, &book); err != nil {
		b.resp.writeJSON(rw, req, http.StatusBadRequest, ErrDecoding)
		b.resp.Errorw("Failed decoding book data from request.", "error", err)

		return
	}

	if err := book.OK(); err != nil {
		b.resp.writeJSON(rw, req, http.StatusBadRequest, err)
		b.resp.Debugw("Failed validating book.", "error", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	if err := b.logic.Import(ctx, book); err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			b.resp.writeJSON(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrRecordNotFound):
			b.resp.writeJSON(rw, req, http.StatusExpectationFailed, fmt.Errorf("author : %w", exceptions.ErrRecordNotFound))
		default:
			b.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		b.resp.Errorw("Failed importing book.", "error", err)

		return
	}

	msg := response{Message: "Book imported successfully."}
	b.resp.writeJSON(rw, req, http.StatusCreated, msg)
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
			b.resp.writeJSON(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrRecordNotFound):
			b.resp.writeJSON(rw, req, http.StatusExpectationFailed, exceptions.ErrRecordNotFound)
		default:
			b.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		b.resp.Errorw("Failed fetching book.", "error", err)

		return
	}

	b.resp.writeJSON(rw, req, http.StatusOK, book)
	b.resp.Debugf("Book fetched successfully.")
}

// FindMany handles bulk fetching books by given OData query.
// If $top number not specified maxOnPage will be used.
// In case number of requested books is more than maxOnPage
// then nextLink will be rendered in response.
func (b Book) FindMany(rw http.ResponseWriter, req *http.Request) {
	filter, err := models.NewDataFilter[models.Book](req.URL)
	if err != nil {
		b.resp.writeJSON(rw, req, http.StatusBadRequest, ErrInvalidQuery)
		b.resp.Errorw("Failed parsing query from request URL.", "error", err)

		return
	}

	maxOnPage, err := strconv.Atoi(os.Getenv("BOOKS_MAX_ON_PAGE"))
	if err != nil {
		b.resp.writeJSON(rw, req, http.StatusBadRequest, exceptions.ErrUnexpected)
		b.resp.Errorw("Failed parsing BOOKS_MAX_ON_PAGE.", "error", err)

		return
	}

	delta := filter.Top - maxOnPage
	filter.Top = maxOnPage

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	books, err := b.logic.FetchMany(ctx, *filter)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			b.resp.writeJSON(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
			// not sure about that.
		case errors.Is(err, exceptions.ErrRecordNotFound):
			b.resp.writeJSON(rw, req, http.StatusBadRequest, exceptions.ErrRecordNotFound)
		default:
			b.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
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

	b.resp.writeJSON(rw, req, http.StatusOK, resp)
	b.resp.Debugf("Books fetched successfully.")
}

func (b Book) AddToFavorites(rw http.ResponseWriter, req *http.Request) {
	var book models.Book
	if err := b.resp.decodeBody(req, &book); err != nil {
		b.resp.writeJSON(rw, req, http.StatusBadRequest, ErrDecoding)
		b.resp.Errorw("Failed decoding book data from request.", "error", err)

		return
	}

	token := retrieveToken[models.AccessToken](req)
	reader := models.Reader{ID: token.ReaderID}

	if book.ID == "" || reader.ID == "" {
		msg := response{Message: "ReaderID and BookID  are required fields."}
		b.resp.writeJSON(rw, req, http.StatusGatewayTimeout, msg)
		b.resp.Debugf(msg.Message)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	if err := b.logic.AddToFavorites(ctx, reader, book); err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			b.resp.writeJSON(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrReaderNotFound):
			b.resp.writeJSON(rw, req, http.StatusBadRequest, exceptions.ErrReaderNotFound)
		case errors.Is(err, exceptions.ErrBookNotFound):
			b.resp.writeJSON(rw, req, http.StatusBadRequest, exceptions.ErrBookNotFound)
		case errors.Is(err, exceptions.ErrRecordExists):
			b.resp.writeJSON(rw, req, http.StatusConflict, exceptions.ErrRecordExists)
		default:
			b.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		b.resp.Errorw("Failed adding book to favorites.", "error", err)

		return
	}

	msg := response{Message: "Book successfully imported fo favorites list."}
	b.resp.writeJSON(rw, req, http.StatusCreated, msg)
	b.resp.Debugf(msg.Message)
}

func (b Book) AddToWishlist(rw http.ResponseWriter, req *http.Request) {
	var book models.Book
	if err := b.resp.decodeBody(req, &book); err != nil {
		b.resp.writeJSON(rw, req, http.StatusBadRequest, ErrDecoding)
		b.resp.Errorw("Failed decoding book data from request.", "error", err)

		return
	}

	token := retrieveToken[models.AccessToken](req)
	reader := models.Reader{ID: token.ReaderID}

	if book.ID == "" || reader.ID == "" {
		msg := response{Message: "ReaderID and BookID  are required fields."}
		b.resp.writeJSON(rw, req, http.StatusGatewayTimeout, msg)
		b.resp.Debugf(msg.Message)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	if err := b.logic.AddToWishlist(ctx, reader, book); err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			b.resp.writeJSON(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrReaderNotFound), errors.Is(err, exceptions.ErrBookNotFound):
			b.resp.writeJSON(rw, req, http.StatusBadRequest, exceptions.ErrRecordNotFound)
		case errors.Is(err, exceptions.ErrRecordExists):
			b.resp.writeJSON(rw, req, http.StatusConflict, exceptions.ErrRecordExists)
		default:
			b.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		b.resp.Errorw("Failed adding book to wishlist.", "error", err)

		return
	}

	msg := response{Message: "Book successfully imported fo wishlist."}
	b.resp.writeJSON(rw, req, http.StatusCreated, msg)
	b.resp.Debugf(msg.Message)
}

func (b Book) ExportToCSV(rw http.ResponseWriter, req *http.Request) {
	filter, err := models.NewDataFilter[models.Book](req.URL)
	if err != nil {
		b.resp.writeJSON(rw, req, http.StatusBadRequest, ErrInvalidQuery)
		b.resp.Errorw("Failed parsing query from request URL.", "error", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	csvArr, err := b.logic.ExportToCSV(ctx, *filter)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			b.resp.writeJSON(rw, req, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		default:
			b.resp.writeJSON(rw, req, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		b.resp.Errorw("Failed exporting books to csv.", "error", err)

		return
	}

	rw.Header().Set("Content-Disposition", "attachment; filename=books.csv")
	rw.Header().Set("Content-Type", "text/csv")
	rw.Header().Set("Transfer-Encoding", "chunked")

	if _, err := rw.Write(csvArr); err != nil {
		b.resp.Errorw("Failed writing response from buffer.", "error", err)

		return
	}

	b.resp.Debugf("Books exported to csv successfully.")
}
