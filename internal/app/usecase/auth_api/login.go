package auth_api

import (
	"auth-api/internal/app/constants"
	localerrors "auth-api/internal/app/errors"
	"auth-api/internal/app/provider/db"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
	token, err := GenerateAccessToken(user.Phone, user.IsEmployee)
	if err != nil {
		return nil, localerrors.NewInternalErr(err)
	}
	refreshToken, err := GenerateRefreshToken(user.Phone, user.IsEmployee)
	if err != nil {
		return nil, localerrors.NewInternalErr(err)
	}
	return &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil

}

type User struct {
	UserId     int
	Phone      string
	Email      string
	Password   string
	IsEmployee bool
	Role       string
}

func (u *authUseCase) searchUser(ctx context.Context, phone string) (*User, error) {
	user, err := u.db.CheckUserExist(ctx, nil, db.CheckUserExistRequest{
		Phone: phone,
	})
	if err != nil {
		return nil, err
	}
	return &User{
		UserId:     user.UserID,
		Phone:      user.Phone,
		Email:      user.Email,
		Password:   user.Password,
		IsEmployee: user.IsEmployee,
		Role:       user.Role,
	}, nil
}

func GenerateAccessToken(phone string, isEmployee bool) (string, error) {
	var jwtSecret = []byte("super-secret-key")
	claims := jwt.MapClaims{
		"phone":       phone,
		"is_employee": isEmployee,
		"exp":         time.Now().Add(15 * time.Minute).Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(phone string, isEmployee bool) (string, error) {
	var jwtSecret = []byte("super-secret-key")
	claims := jwt.MapClaims{
		"phone":       phone,
		"is_employee": isEmployee,
		"exp":         time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
