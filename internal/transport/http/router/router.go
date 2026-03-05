package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/handler"
	customMiddleware "github.com/mamochiro/go-microservice-template/internal/transport/http/middleware"
)

func NewRouter(healthHandler *handler.HealthHandler, userHandler *handler.UserHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.CorrelationID)

	r.Get("/health", healthHandler.Check)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.Create)
			r.Get("/", userHandler.List)
			r.Get("/{id}", userHandler.Get)
			r.Put("/{id}", userHandler.Update)
			r.Delete("/{id}", userHandler.Delete)
		})
	})

	return r
}
