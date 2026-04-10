package auth_api

import (
	pgclient "auth-api/internal/app/client/pg"
	"auth-api/internal/app/constants"
	localerrors "auth-api/internal/app/errors"
	"auth-api/internal/app/provider/db"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	Phone      string
	Email      string
	Password   string
	IsEmployee bool
}

func (u *authUseCase) Register(ctx context.Context, req Register) localerrors.Error {
	checkUserExist, err := u.db.CheckUserExist(ctx, nil, db.CheckUserExistRequest{
		Phone: req.Phone,
	})
	if err != nil {
		return localerrors.NewInternalErr(err)
	}
	if checkUserExist != nil {
		return localerrors.NewInternalErr(errors.New(constants.UserAlreadyExistsError))
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return localerrors.NewInternalErr(err)
	}

	txErr := u.db.WithTransaction(ctx, func(ctx context.Context, tx pgclient.Transaction) error {
		role := "client"
		if req.IsEmployee == true {
			employeeRole, err := u.db.GetUserRoleByPhone(ctx, tx, db.GetUserRoleRequest{
				Phone: req.Phone,
			})
			if err != nil {
				return err
			}

			if employeeRole != nil {
				role = employeeRole.Role
			} else {
				return errors.New(constants.EmployeeNotFoundError)
			}
		}

		if err := u.db.CreateUser(ctx, tx, db.CreateUserRequest{
			Phone:      req.Phone,
			Email:      req.Email,
			Password:   string(hashedPassword),
			IsEmployee: req.IsEmployee,
			Role:       role,
		}); err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		return localerrors.NewInternalErr(txErr)
	}

	return nil
}
