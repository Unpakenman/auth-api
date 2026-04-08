package db

import (
	pgclient "auth-api/internal/app/client/pg"
	"context"
)

type authDBProvider struct {
	conn pgclient.PGClient
}

func NewAuthDBProvider(dbConn pgclient.PGClient) AuthProvider {
	return &authDBProvider{
		conn: dbConn,
	}
}

func (p *authDBProvider) WithTransaction(ctx context.Context, fn func(context.Context, pgclient.Transaction) error) error {
	return p.conn.WithTransaction(ctx, fn)
}
