//go:build wireinject
// +build wireinject

package app

import (
	"net/http"

	"github.com/google/wire"
	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/mamochiro/go-microservice-template/internal/domain/service"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/cache"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/database"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/repository"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/handler"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/router"
)

func InitializeApp(cfg *config.Config) (http.Handler, func(), error) {
	wire.Build(
		database.NewPostgresDB,
		cache.NewRedisClient,
		cache.NewCacheRepository,
		repository.NewUserRepository,
		service.NewUserService,
		handler.NewHealthHandler,
		handler.NewUserHandler,
		router.NewRouter,
	)
	return nil, nil, nil
}
