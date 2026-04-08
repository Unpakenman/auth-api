package bootstrap

import (
	"auth-api/internal/app/client/pg"
	"auth-api/internal/app/client/redis"
	"auth-api/internal/app/config"
	"auth-api/internal/app/grpcserver"
	"auth-api/internal/app/grpcserver/mapper"
	logger "auth-api/internal/app/log"
	cacheprovider "auth-api/internal/app/provider/cache/redis"
	"auth-api/internal/app/provider/db"
	"auth-api/internal/app/usecase/auth_api"
	"auth-api/internal/app/validator"
	"context"
	"fmt"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"os"
	"os/signal"
	"syscall"
)

func RunService(ctx context.Context, cfg *config.Values, log logger.LogClient) {

	exit := make(chan os.Signal, 1)

	mapperInstance := mapper.New()
	validatorInstance := validator.New()
	dbConn, err := pg.New(cfg.ClinicsDB)
	if err != nil {
		log.Error(err)
	}

	authDBProvider := db.NewAuthDBProvider(dbConn)

	cacheClient, err := redis.NewRedisClient(cfg.Redis.URL)
	if err != nil {
		log.Fatal(err)
	}
	cacheProvider := cacheprovider.NewRedisCache(cacheClient, cfg.Redis.Prefix)

	authUseCaseInstance := auth_api.NewUseCase(authDBProvider, log, cfg, cacheProvider)

	grpcPortListener, err := NewGRPCPortListener(cfg.GRPCServer)
	if err != nil {
		log.Error(err)
	}
	defer func() {
		err := grpcPortListener.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	authServerInstance := grpcserver.NewAuthServer(
		log,
		validatorInstance,
		mapperInstance,
		authUseCaseInstance)
	healthcheck := health.NewServer()
	grpcServer, err := NewGRPCServer(cfg.GRPCServer, log)
	if err != nil {
		log.Error(err)
	}

	pb.RegisterAuthServer(grpcServer, authServerInstance)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthcheck)
	reflection.Register(grpcServer)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(grpcPortListener); err != nil {
			log.Error(err, "grpc serve failed")
		}
	}()

	log.Info("app service started")
	select {
	case v := <-exit:
		log.Warn(fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		log.InfoCtx(ctx, "ctx.Done: ", done)
	}
	grpcServer.GracefulStop()

	if err := dbConn.CloseConnections(); err != nil {
		log.Error(err, "failed to close database connection")
	}
	log.Info("Server Exited Properly")
}
