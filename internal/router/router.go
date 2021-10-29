package router

import (
	//"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		//r.GET("/rate",  )
	})

	return r
}
