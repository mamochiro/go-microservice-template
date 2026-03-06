//go:build wireinject
// +build wireinject

package app

import (
	"net/http"

	"github.com/google/wire"
	"github.com/mamochiro/go-microservice-template/internal/config"
	"github.com/mamochiro/go-microservice-template/internal/domain/service"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/repository"
	transport "github.com/mamochiro/go-microservice-template/internal/transport/http"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/router"
)

func InitializeApp(cfg *config.Config) (http.Handler, func(), error) {
	wire.Build(
		infrastructure.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		transport.ProviderSet,
		router.NewRouter,
	)
	return nil, nil, nil
}
