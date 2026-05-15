package auth_api

import (
	"auth-api/internal/app/constants"
	appErrors "auth-api/internal/app/errors"
	localerrors "auth-api/internal/app/errors"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
)

type LogoutRequest struct {
	Token        string
	RefreshToken string
}

func (u authUseCase) Logout(ctx context.Context, req LogoutRequest,
) localerrors.Error {
	signatureToken := constants.PathToAccessSignature
	signatureRefreshToken := constants.PathToRefreshTokenSignature

	accessClaims, err := parseToken(req.Token, signatureToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return localerrors.NewInternalErr(appErrors.InvalidSignature)
		default:
			return localerrors.NewInternalErr(err)
		}
	}

	refreshClaims, err := parseToken(req.RefreshToken, signatureRefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrSignatureInvalid):

			return localerrors.NewInternalErr(appErrors.InvalidSignature)
		default:
			return localerrors.NewInternalErr(err)
		}
	}

	if accessClaims.UserId != refreshClaims.UserId {
		return localerrors.NewInternalErr(appErrors.InvalidToken)
	}

	cacheKey := "refresh_token" + strconv.FormatInt(accessClaims.UserId, 10)
	if err := u.cache.Delete(ctx, cacheKey); err != nil {
		return localerrors.NewInternalErr(err)
	}

	return nil
}
