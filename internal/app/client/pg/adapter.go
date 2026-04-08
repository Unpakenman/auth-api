package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
)

const (
	LogLevelTrace = 6
	LogLevelDebug = 5
	LogLevelInfo  = 4
	LogLevelWarn  = 3
	LogLevelError = 2
	LogLevelNone  = 1
)

type LogLevel int

// LogAdapter adapter needed for use in pgx
type LogAdapter struct {
	log LogClient
}

type LogClient interface {
	InfoCtx(ctx context.Context, msg string, fields ...interface{})
	TraceCtx(ctx context.Context, msg string, fields ...interface{})
	WarnCtx(ctx context.Context, msg string, fields ...interface{})
	ErrorCtx(ctx context.Context, err error, fields ...interface{})
	DebugCtx(ctx context.Context, msg string, fields ...interface{})
}

func NewLogAdapter(l LogClient) *LogAdapter {
	return &LogAdapter{log: l}
}

// Log common.go log method
func (l *LogAdapter) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {

	switch level {
	case pgx.LogLevelTrace:
		l.log.TraceCtx(ctx, msg)
	case pgx.LogLevelDebug:
		l.log.DebugCtx(ctx, msg)
	case pgx.LogLevelInfo:
		l.log.InfoCtx(ctx, msg)
	case pgx.LogLevelWarn:
		l.log.WarnCtx(ctx, msg)
	case pgx.LogLevelError:
		l.log.ErrorCtx(ctx, errors.New(msg))
	default:
		l.log.ErrorCtx(ctx, fmt.Errorf("INVALID_PGX_LOG_LEVEL, %d", level))
	}
}
