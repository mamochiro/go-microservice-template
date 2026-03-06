package http

import (
	"github.com/google/wire"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/auth"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/health"
	"github.com/mamochiro/go-microservice-template/internal/transport/http/user"
)

// ProviderSet is http handler providers.
var ProviderSet = wire.NewSet(
	health.NewHandler,
	user.NewHandler,
	auth.NewHandler,
)
