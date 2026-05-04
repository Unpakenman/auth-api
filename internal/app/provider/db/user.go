package db

import (
	pgclient "auth-api/internal/app/client/pg"
	"auth-api/internal/app/provider/db/models"
	"context"
	"database/sql"
	"errors"
)

type CheckUserExistRequest struct {
	Phone string
}

func (p *authDBProvider) CheckUserExist(
	ctx context.Context,
	tx pgclient.Transaction,
	req CheckUserExistRequest,
) (*models.UserExistResponse, error) {
	var userExist models.UserExistResponse
	err := p.conn.NamedGetContext(
		ctx,
		&userExist,
		"CheckUserExist",
		nil,
		tx,
		req.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &userExist, nil
}

type GetUserRoleRequest struct {
	Phone string
}

func (p *authDBProvider) GetUserRoleByPhone(
	ctx context.Context,
	tx pgclient.Transaction,
	req GetUserRoleRequest,
) (*models.EmployeeRoleResponse, error) {
	var role models.EmployeeRoleResponse
	err := p.conn.NamedGetContext(
		ctx,
		&role,
		"GetEmployeeRoleByPhone",
		nil,
		tx,
		req.Phone)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

type CreateUserRequest struct {
	Phone      string
	Email      string
	Password   string
	IsEmployee bool
	Role       string
}

func (p *authDBProvider) CreateUser(
	ctx context.Context,
	tx pgclient.Transaction,
	req CreateUserRequest,
) error {
	_, err := p.conn.Exec(
		ctx,
		"CreateUser",
		nil,
		tx,
		req.Phone,
		req.Email,
		req.Password,
		req.IsEmployee,
		req.Role)
	if err != nil {
		return err
	}
	return nil
}
