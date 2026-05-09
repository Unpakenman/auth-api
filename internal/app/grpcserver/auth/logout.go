package auth

import (
	localerrors "auth-api/internal/app/errors"
	"context"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
)

func (s *ServerAuth) Logout(ctx context.Context, req *pb.LogoutRequest,
) (*pb.LogoutResponse, error) {
	if errs := s.validator.Logout(req); errs != nil {
		err := localerrors.NewInvalidArgumentErr(*errs)
		s.log.InfoCtx(ctx, "validator logout error", err.Error())
		return nil, s.mapper.ResultErrorToProtoError(err)
	}

	useCaseReq := s.mapper.LogoutToUseCase(req)
	if err := s.authUseCase.Logout(ctx, useCaseReq); err != nil {
		return nil, s.mapper.ResultErrorToProtoError(err)
	}

	return &pb.LogoutResponse{}, nil
}
