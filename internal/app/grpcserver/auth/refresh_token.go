package auth

import (
	localerrors "auth-api/internal/app/errors"
	"context"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
)

func (s *ServerAuth) RefreshToken(ctx context.Context, req *pb.RefreshAccessTokenRequest,
) (*pb.RefreshAccessTokenResponse, error) {
	if errs := s.validator.RefreshAccessToken(req); errs != nil {
		err := localerrors.NewInvalidArgumentErr(*errs)
		s.log.InfoCtx(ctx, "validator refresh_token error", err.Error())
		return nil, s.mapper.ResultErrorToProtoError(err)
	}

	useCaseReq := s.mapper.RefreshAccessTokenToUseCase(req)
	useCaseResp, err := s.authUseCase.RefreshAccessToken(ctx, useCaseReq)
	if err != nil {
		errs := localerrors.NewInternalErr(err)
		s.log.InfoCtx(ctx, "auth refresh_token error", err.Error())
		return nil, s.mapper.ResultErrorToProtoError(errs)
	}

	return &pb.RefreshAccessTokenResponse{
		Token:        useCaseResp.Token,
		RefreshToken: useCaseResp.RefreshToken,
	}, nil
}
