package auth_api

import (
	"auth-api/internal/app/config"
	localerrors "auth-api/internal/app/errors"
	logger "auth-api/internal/app/log"
	providerCache "auth-api/internal/app/provider/cache/redis"
	"auth-api/internal/app/provider/db"
	"context"
)

type authUseCase struct {
	db     db.AuthProvider
	logger logger.LogClient
	config *config.Values
	cache  providerCache.Cache
}

func NewUseCase(
	provider db.AuthProvider,
	logger logger.LogClient,
	config *config.Values,
	providerCache providerCache.Cache,
) UseCase {
	return &authUseCase{
		db:     provider,
		logger: logger,
		config: config,
		cache:  providerCache,
	}
}

type UseCase interface {
	Register(ctx context.Context, req Register) localerrors.Error
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, localerrors.Error)
}
