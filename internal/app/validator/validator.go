package validator

import (
	localerrors "auth-api/internal/app/errors"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
)

type validator struct{}

type Validator interface {
	Register(req *pb.RegisterRequest) *[]localerrors.FieldViolation
	Login(req *pb.LoginRequest) *[]localerrors.FieldViolation
}

func New() Validator {
	return &validator{}
}
