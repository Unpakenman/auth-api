package main

import (
	"auth-api/internal/app/config"
	appLog "auth-api/internal/app/log"
	"auth-api/internal/bootstrap"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	logger, err := appLog.New(*cfg)
	config.Config = cfg
	bootstrap.RunService(ctx, cfg, logger)
}
