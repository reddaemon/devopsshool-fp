package router

import (
	"final-project/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRouter(handler *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	r.Route("/v1", func(r chi.Router) {
		r.Get("/rate", handler.GetRate)
		r.Get("/pullrate/{dd:\\d\\d}/{mm:\\d\\d}/{yyyy:\\d\\d\\d\\d}", handler.PullRate)
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
		r.Post("/logout", handler.Logout)
	})

	return r
}
