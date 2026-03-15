package infrastructure

import (
	"github.com/google/wire"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/cache"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/database"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/email"
)

// ProviderSet is infrastructure providers.
var ProviderSet = wire.NewSet(
	database.NewPostgresDB,
	cache.NewRedisClient,
	email.NewResendService,
)
