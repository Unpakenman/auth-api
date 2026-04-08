package pg

import (
	"auth-api/internal/app/config"
	"embed"
	"fmt"
)

var (
	//go:embed queries/*
	queryFiles embed.FS

	pathsToDbQueries = []string{"queries/"}
)

//go:generate ../../../../bin/mockery --srcpkg=vcs.bingo-boom.ru/digital_department/ru/application-layer/resources/dependencies/go-modules/pgclient --case=underscore --name=Transaction

func New(cfg *config.DBConfig) (PGClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("invalid pg config")
	}
	connString := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v",
		cfg.User,
		cfg.Password,
		cfg.Hostname,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
	configValues := PostgreSQL{
		ConnString:                      connString,
		PathsToQueries:                  pathsToDbQueries,
		MaxOpenConns:                    cfg.MaxOpenConns,
		MaxIdleConns:                    cfg.MaxIdleConns,
		MaxLifeTimeConns:                cfg.MaxLifeTimeConns,
		StatementTimeout:                cfg.StatementTimeout,
		IdleInTransactionSessionTimeout: cfg.IdleInTransactionSessionTimeout,
		LockTimeout:                     cfg.LockTimeout,
		LogLevel:                        LogLevelNone,
	}
	return NewClient(configValues, queryFiles, nil)
}
