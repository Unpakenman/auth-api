package db

import (
	pgclient "auth-api/internal/app/client/pg"
	"auth-api/internal/app/provider/db/models"
	"context"
)

//go:generate ../../../../bin/mockery --with-expecter --case=underscore --name=GoExampleProvider

type AuthProvider interface {
	WithTransaction(ctx context.Context, fn func(context.Context, pgclient.Transaction) error) error
	CheckUserExist(
		ctx context.Context,
		tx pgclient.Transaction,
		req CheckUserExistRequest,
	) (*models.UserExistResponse, error)
	GetUserRoleByPhone(
		ctx context.Context,
		tx pgclient.Transaction,
		req GetUserRoleRequest,
	) (*models.EmployeeRoleResponse, error)
	CreateUser(
		ctx context.Context,
		tx pgclient.Transaction,
		req CreateUserRequest,
	) error
}
