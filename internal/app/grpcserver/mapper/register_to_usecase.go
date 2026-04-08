package mapper

import (
	"auth-api/internal/app/usecase/auth_api"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
)

func (m *mapper) RegisterToUsecase(req *pb.RegisterRequest) auth_api.Register {
	return auth_api.Register{
		Phone:      req.PhoneNumber,
		Email:      req.Email,
		Password:   req.Password,
		IsEmployee: req.IsAmployee,
	}
}
