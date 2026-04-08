package db

import (
	pgclient "auth-api/internal/app/client/pg"
	"context"
)

//go:generate ../../../../bin/mockery --with-expecter --case=underscore --name=GoExampleProvider

type AuthProvider interface {
	WithTransaction(ctx context.Context, fn func(context.Context, pgclient.Transaction) error) error
}
