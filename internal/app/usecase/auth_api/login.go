package auth_api

import (
	"auth-api/internal/app/constants"
	localerrors "auth-api/internal/app/errors"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type LoginRequest struct {
	Phone      string
	Password   string
	IsEmployee bool
}
type LoginResponse struct {
	Token        string
	RefreshToken string
}

func (u *authUseCase) Login(ctx context.Context, req LoginRequest,
) (*LoginResponse, localerrors.Error) {
	user, err := u.searchUser(ctx, req.Phone)
	if err != nil {
		return nil, localerrors.NewBadRequestErr(err)
	}

	if user == nil {
		return nil, localerrors.NewBadRequestErr(errors.New(constants.UserNotFoundError))
	}

	passwordsMatchError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if passwordsMatchError != nil {
		return nil, localerrors.NewBadRequestErr(errors.New(constants.PasswordsMismatch))
	}

	token, err := GenerateAccessToken(user.UserId)
	if err != nil {
		return nil, localerrors.NewInternalErr(err)
	}

	cacheKey := "refresh_token:" + strconv.FormatInt(user.UserId, 10)
	refreshToken, err := GenerateRefreshToken(user.UserId)
	if err != nil {
		return nil, localerrors.NewInternalErr(err)
	}

	if err := u.cache.Delete(ctx, cacheKey); err != nil {
		return nil, localerrors.NewInternalErr(err)
	} //либо не удалять чтобы не было разлогина на других устройствах
	if err := u.cache.Set(ctx, cacheKey, []byte(refreshToken), time.Hour*2160); err != nil {
		return nil, localerrors.NewInternalErr(err)
	}

	return &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}
