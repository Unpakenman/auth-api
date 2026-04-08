package auth

import (
	"auth-api/internal/app/grpcserver/mapper"
	logger "auth-api/internal/app/log"
	"auth-api/internal/app/usecase/auth_api"
	"auth-api/internal/app/validator"
)

type ServerAuth struct {
	log         logger.LogClient
	validator   validator.Validator
	mapper      mapper.Mapper
	authUseCase auth_api.UseCase
}

func NewServer(
	logger logger.LogClient,
	validator validator.Validator,
	mapper mapper.Mapper,
	authUseCase auth_api.UseCase,
) *ServerAuth {
	return &ServerAuth{
		log:         logger,
		validator:   validator,
		mapper:      mapper,
		authUseCase: authUseCase,
	}
}
