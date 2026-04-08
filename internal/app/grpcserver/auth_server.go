package grpcserver

import (
	"auth-api/internal/app/grpcserver/auth"
	"auth-api/internal/app/grpcserver/mapper"
	logger "auth-api/internal/app/log"
	"auth-api/internal/app/usecase/auth_api"
	"auth-api/internal/app/validator"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth"
)

type AuthServer struct {
	*auth.ServerAuth
}

func NewAuthServer(
	logger logger.LogClient,
	validator validator.Validator,
	mapper mapper.Mapper,
	clinicUseCase auth_api.UseCase,
) pb.AuthServer {
	return &AuthServer{
		auth.NewServer(logger, validator, mapper, clinicUseCase),
	}
}
