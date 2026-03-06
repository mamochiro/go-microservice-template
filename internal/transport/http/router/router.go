package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	_ "github.com/mamochiro/go-microservice-template/docs"
	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/handler"
	customMiddleware "github.com/mamochiro/go-microservice-template/internal/transport/http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewRouter(cfg *config.Config, healthHandler *handler.HealthHandler, userHandler *handler.UserHandler, authHandler *handler.AuthHandler) http.Handler {
	r := chi.NewRouter()

	// Basic Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Security: Secure Headers
	r.Use(customMiddleware.SecureHeaders(cfg.App.Env))

	// Rate Limiting
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// Security: CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Correlation-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Custom Middlewares
	r.Use(customMiddleware.CorrelationID)
	r.Use(customMiddleware.Metrics)
	r.Use(customMiddleware.Logger)

	// Public Routes
	r.Get("/health", healthHandler.Check)
	r.Handle("/metrics", promhttp.Handler())
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/signup", userHandler.Create) // Use /signup for public creation

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.Auth(cfg.App.JWTSecret))

			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.List)
				r.Get("/{id}", userHandler.Get)
				r.Put("/{id}", userHandler.Update)
				r.Delete("/{id}", userHandler.Delete)
			})
		})
	})

	return otelhttp.NewHandler(r, "http-server")
}
