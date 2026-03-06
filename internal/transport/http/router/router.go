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
	"github.com/mamochiro/go-microservice-template/internal/transport/http/auth"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/health"
	mw "github.com/mamochiro/go-microservice-template/internal/transport/http/middleware"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/user"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewRouter(cfg *config.Config, healthHandler *health.Handler, userHandler *user.Handler, authHandler *auth.Handler) http.Handler {
	r := chi.NewRouter()

	// 1. Global Middlewares
	setupGlobalMiddlewares(r, cfg)

	// 2. Public System Routes
	r.Get("/health", healthHandler.Check)
	r.Handle("/metrics", promhttp.Handler())

	// 3. Swagger Documentation
	setupSwaggerRoutes(r)

	// 4. API v1 Routes
	r.Route("/api/v1", func(r chi.Router) {
		auth.Register(r, authHandler)
		user.Register(r, userHandler, cfg.App.JWTSecret)
	})

	return otelhttp.NewHandler(r, "http-server")
}

func setupGlobalMiddlewares(r *chi.Mux, cfg *config.Config) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Security: Secure Headers
	r.Use(mw.SecureHeaders(cfg.App.Env))

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
	r.Use(mw.CorrelationID)
	r.Use(mw.Metrics)
	r.Use(mw.Logger)
}

func setupSwaggerRoutes(r chi.Router) {
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)
}
