package router

import (
	"final-project/internal/handlers"
	"final-project/internal/middleware"
	"net/http"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterRouter(handler *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(chiprometheus.NewMiddleware("my-api"))

	//r.Use(middleware.Logger)
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	r.Handle("/metrics", promhttp.Handler())
	r.Route("/v1", func(r chi.Router) {
		dataGroup := r.Group(nil)
		dataGroup.Use(middleware.TokenAuthMiddleware)
		dataGroup.Get("/rate", (handler.GetRate))
		dataGroup.Get("/pullrate/{dd:\\d\\d}/{mm:\\d\\d}/{yyyy:\\d\\d\\d\\d}", handler.PullRate)
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
		r.Post("/logout", handler.Logout)
		r.Post("/refresh", handler.Refresh)
	})

	return r
}
