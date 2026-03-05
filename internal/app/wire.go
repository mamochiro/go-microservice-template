//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/mamochiro/go-microservice-template/internal/domain/service"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/cache"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/database"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/repository"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/handler"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/router"
	"github.com/go-chi/chi/v5"
)

func InitializeApp(cfg *config.Config) (*chi.Mux, func(), error) {
	wire.Build(
		database.NewPostgresDB,
		cache.NewRedisClient,
		repository.NewUserRepository,
		service.NewUserService,
		handler.NewHealthHandler,
		handler.NewUserHandler,
		router.NewRouter,
	)
	return nil, nil, nil
}
