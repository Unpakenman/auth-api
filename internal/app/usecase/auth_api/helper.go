package auth_api

import (
	"auth-api/internal/app/constants"
	"auth-api/internal/app/provider/db"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
	"time"
)

type User struct {
	UserId     int64
	Phone      string
	Email      string
	Password   string
	IsEmployee bool
	Role       string
}

func (u *authUseCase) searchUser(ctx context.Context, phone string) (*User, error) {
	phone = strings.NewReplacer(" ", "", "-", "", "+", "").Replace(phone)
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

func GenerateAccessToken(userId int64) (string, error) {
	var jwtSecret = []byte(os.Getenv(constants.PathToAccessSignature))
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(userId int64) (string, error) {
	var jwtSecret = []byte(os.Getenv(constants.PathToRefreshTokenSignature))
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

type TokenClaims struct {
	UserId int64 `json:"user_id"`
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
