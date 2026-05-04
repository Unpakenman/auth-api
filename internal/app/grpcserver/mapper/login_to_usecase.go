package mapper

import (
	"auth-api/internal/app/usecase/auth_api"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
)

func (m *mapper) LoginToUseCase(request *pb.LoginRequest) auth_api.LoginRequest {
	return auth_api.LoginRequest{
		Phone:      request.PhoneNumber,
		Password:   request.Password,
		IsEmployee: request.IsAmployee,
	}
}
