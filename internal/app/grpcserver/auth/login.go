package auth

import (
	localerrors "auth-api/internal/app/errors"
	"context"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
)

func (s *ServerAuth) Login(ctx context.Context, req *pb.LoginRequest,
) (*pb.LoginResponse, error) {
	if errs := s.validator.Login(req); errs != nil {
		err := localerrors.NewInvalidArgumentErr(*errs)
		s.log.InfoCtx(ctx, "validator login error", err.Error())
		return nil, s.mapper.ResultErrorToProtoError(err)
	}
	useCaseReq := s.mapper.LoginToUseCase(req)
	usecaseResp, err := s.authUseCase.Login(ctx, useCaseReq)
	if err != nil {
		errs := localerrors.NewInternalErr(err)
		s.log.InfoCtx(ctx, "auth usecase login error", err.Error())
		return nil, s.mapper.ResultErrorToProtoError(errs)
	}
	return &pb.LoginResponse{
		Token:        usecaseResp.Token,
		RefreshToken: usecaseResp.RefreshToken,
	}, nil
}
