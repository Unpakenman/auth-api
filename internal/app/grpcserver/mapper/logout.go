package mapper

import (
	"auth-api/internal/app/usecase/auth_api"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
)

func (m *mapper) LogoutToUseCase(request *pb.LogoutRequest) auth_api.LogoutRequest {
	return auth_api.LogoutRequest{
		Token:        request.Token,
		RefreshToken: request.RefreshToken,
	}
}
