package router

import (
	//"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
	"final-project/internal/handlers"
)

func RegisterRouter(handler *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Get("/rate", handler.GetRate)
	})

	return r
}
