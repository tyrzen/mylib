package rest

import "github.com/go-chi/chi/v5"

func Auth(router chi.Router) {
	router.Route("/auth", func(router chi.Router) {
		router.Post("/login", nil)
		router.Post("/logout", nil)
		router.Post("/token", nil)
	})
}
