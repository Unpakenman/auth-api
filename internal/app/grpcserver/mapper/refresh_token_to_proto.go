package mapper

import (
	"auth-api/internal/app/usecase/auth_api"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
)

func (m *mapper) RefreshAccessTokenToUseCase(request *pb.RefreshAccessTokenRequest) auth_api.RefreshAccessTokenRequest {
	return auth_api.RefreshAccessTokenRequest{
		Token:        request.Token,
		RefreshToken: request.RefreshToken,
	}
}
