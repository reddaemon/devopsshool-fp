package router

import (
	//"github.com/go-chi/chi/v5/middleware"
	"final-project/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRouter(handler *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	r.Route("/v1", func(r chi.Router) {
		r.Get("/rate", handler.GetRate)
	})

	return r
}
