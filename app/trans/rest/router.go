package rest

import (
	"github.com/go-chi/chi/v5"
)

type Router struct{ chi.Router }

func NewMux(routes ...func(chi.Router)) Router {
	rtr := chi.NewRouter()

	for _, r := range routes {
		rtr.Group(r)
	}

	return Router{rtr}
}
