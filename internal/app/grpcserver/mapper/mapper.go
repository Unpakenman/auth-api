package mapper

import (
	localerrors "auth-api/internal/app/errors"
	"auth-api/internal/app/usecase/auth_api"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
)

type mapper struct{}

type Mapper interface {
	RegisterToUsecase(req *pb.RegisterRequest) auth_api.Register
	LoginToUseCase(request *pb.LoginRequest) auth_api.LoginRequest
	ResultErrorToProtoError(resultError localerrors.Error) error
}

func New() Mapper { return &mapper{} }
