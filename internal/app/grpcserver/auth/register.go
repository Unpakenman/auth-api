package auth

import (
	localerrors "auth-api/internal/app/errors"
	"context"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
)

func (s *ServerAuth) Register(ctx context.Context, req *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {
	if errs := s.validator.Register(req); errs != nil {
		err := localerrors.NewInvalidArgumentErr(*errs)
		s.log.InfoCtx(ctx, "validator register error", err.Error())
		return nil, s.mapper.ResultErrorToProtoError(err)
	}
	useCaseReq := s.mapper.RegisterToUsecase(req)
	if err := s.authUseCase.Register(ctx, useCaseReq); err != nil {
		s.log.ErrorCtx(ctx, err)
		return nil, s.mapper.ResultErrorToProtoError(err)
	}
	return &pb.RegisterResponse{}, nil
}
