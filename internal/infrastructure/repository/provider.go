package repository

import (
	"github.com/google/wire"
	"github.com/mamochiro/go-microservice-template/internal/infrastructure/cache"
)

// ProviderSet is repository providers.
var ProviderSet = wire.NewSet(
	cache.NewCacheRepository,
	NewUserRepository,
	NewCachedUserRepository,
)
