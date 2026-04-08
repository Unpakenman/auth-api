package auth_api

import (
	localerrors "auth-api/internal/app/errors"
	"context"
)

type Register struct {
	Phone      string
	Email      string
	Password   string
	IsEmployee bool
}

func (u *authUseCase) Register(ctx context.Context, req Register) localerrors.Error {
	return nil
}
