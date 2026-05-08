package auth_api

import (
	"auth-api/internal/app/constants"
	"auth-api/internal/app/provider/db"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

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
	var jwtSecret = []byte(os.Getenv(constants.PathToAccessSignature))
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
	var jwtSecret = []byte(os.Getenv(constants.PathToRefreshTokenSignature))
	claims := jwt.MapClaims{
		"phone":       phone,
		"is_employee": isEmployee,
		"exp":         time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

type TokenClaims struct {
	Phone      string `json:"phone"`
	IsEmployee bool   `json:"is_employee"`
	jwt.RegisteredClaims
}

func parseToken(tokenString string, signaturePath string) (*TokenClaims, error) {
	secret := []byte(os.Getenv(signaturePath))

	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
