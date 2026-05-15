package auth_api

import (
	"auth-api/internal/app/constants"
	appErrors "auth-api/internal/app/errors"
	localerrors "auth-api/internal/app/errors"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type RefreshAccessTokenRequest struct {
	Token        string
	RefreshToken string
}

type RefreshAccessTokenResponse struct {
	Token        string
	RefreshToken string
}

func (u *authUseCase) RefreshAccessToken(ctx context.Context, req RefreshAccessTokenRequest,
) (*RefreshAccessTokenResponse, localerrors.Error) {
	signatureToken := constants.PathToAccessSignature
	signatureRefreshToken := constants.PathToRefreshTokenSignature

	accessClaims, err := parseToken(req.Token, signatureToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return &RefreshAccessTokenResponse{}, localerrors.NewInternalErr(appErrors.InvalidSignature)
		default:
			return &RefreshAccessTokenResponse{}, localerrors.NewInternalErr(err)
		}
	}

	refreshClaims, err := parseToken(req.RefreshToken, signatureRefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return &RefreshAccessTokenResponse{}, localerrors.NewInternalErr(appErrors.InvalidSignature)
		case errors.Is(err, jwt.ErrTokenExpired):
			return &RefreshAccessTokenResponse{}, localerrors.NewInternalErr(appErrors.TokenExpired)
		default:
			return &RefreshAccessTokenResponse{}, localerrors.NewInternalErr(err)
		}
	}

	if accessClaims.UserId != refreshClaims.UserId {
		return &RefreshAccessTokenResponse{}, localerrors.NewInternalErr(appErrors.InvalidToken)
	}

	token, err := GenerateAccessToken(accessClaims.UserId)
	if err != nil {
		return nil, localerrors.NewInternalErr(err)
	}

	cacheKey := "refresh_token" + strconv.FormatInt(accessClaims.UserId, 10)
	_, err = u.cache.Get(ctx, cacheKey)
	if err != nil {
		return &RefreshAccessTokenResponse{}, localerrors.NewInternalErr(appErrors.TokenNotFound)
	}

	if err := u.cache.Delete(ctx, cacheKey); err != nil {
		return nil, localerrors.NewInternalErr(err)
	}

	refreshToken, err := GenerateRefreshToken(refreshClaims.UserId)
	if err != nil {
		return nil, localerrors.NewInternalErr(err)
	}

	if err := u.cache.Set(ctx, cacheKey, []byte(refreshToken), time.Hour*2160); err != nil {
		return nil, localerrors.NewInternalErr(err)
	}

	return &RefreshAccessTokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}
