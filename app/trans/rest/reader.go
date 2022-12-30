package rest

import (
	"fmt"
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

func (r *Reader) Route(router chi.Router) {
	router.Get("/", r.handleDummy())
}

func (r *Reader) handleDummy() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "bla bla lba")
	}
}
