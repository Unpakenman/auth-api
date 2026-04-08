package bootstrap

import (
	"auth-api/internal/app/config"
	localerrors "auth-api/internal/app/errors"
	logger "auth-api/internal/app/log"
	logoptions "auth-api/internal/app/log/options"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

func NewGRPCPortListener(cfg *config.GRPCServerConfig,
) (net.Listener, error) {
	if cfg == nil {
		return nil, fmt.Errorf("invalid grpc config")
	}
	return net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
}

func NewGRPCServer(
	cfg *config.GRPCServerConfig,
	log logger.LogClient,
	interceptors ...grpc.ServerOption,
) (*grpc.Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("invalid grpc config")
	}
	interceptors = append(interceptors, grpc.ChainUnaryInterceptor(SetRequestIDMW(log)))
	interceptors = append(interceptors, grpc.ChainUnaryInterceptor(LoggingAndTracingMW(log)))
	interceptors = append(interceptors, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: cfg.KeepaliveMaxConnectionIdle,
		Timeout:           cfg.KeepaliveTimeout,
		MaxConnectionAge:  cfg.KeepaliveMaxConnectionAge,
		Time:              cfg.KeepaliveTime,
	}))
	grpcServer := grpc.NewServer(interceptors...)
	return grpcServer, nil
}

type GRPCResponse[T any, K any] struct {
	Data  *T
	Error *GRPCResponseError[K]
}
type GRPCResponseError[T any] struct {
	Code    int32
	Message string
	Details T
}

func SetRequestIDMW(
	logger logger.LogClient,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx = logger.SetOptionsToCtx(
			ctx,
			logoptions.WithRequestID(uuid.New().String()),
		)
		return handler(ctx, req)
	}
}

func LoggingAndTracingMW(
	logger logger.LogClient,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		defer localerrors.RecoverHandler(ctx, logger)()
		methodName := "invalid name"
		if info != nil {
			methodName = info.FullMethod
		} else {
			logger.ErrorMessage("invalid grpc method name", req)
		}

		span := trace.SpanFromContext(ctx)
		traceID := span.SpanContext().TraceID().String()

		ctx = logger.SetOptionsToCtx(
			ctx,
			logoptions.WithTraceID(traceID),
			logoptions.WithProtocol(logoptions.GRPCProtocol),
		)

		logger.InfoCtx(ctx, "grpc_server.request method: "+methodName, req)

		start := time.Now()
		reply, errReply := handler(ctx, req)
		duration := time.Since(start)

		var code codes.Code
		if errReply != nil {
			code = status.Code(errReply)
		} else {
			code = codes.OK
		}

		requestInfoText := fmt.Sprintf(
			"grpc_server.response method: %s, time spent: %v, code: %s",
			methodName,
			duration,
			code.String(),
		)

		logger.InfoCtx(ctx, requestInfoText, reply)

		return reply, errReply
	}
}
