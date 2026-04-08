package log

import (
	"auth-api/internal/app/config"
	"auth-api/internal/app/log/options"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func New(cfg config.Values) (LogClient, error) {
	log, err := NewLog(Options{
		ConsoleOptions: ConsoleOptions{
			Level: cfg.LogLevel,
		},
	})
	if err != nil {
		return nil, err
	}

	return log, err
}

type Options struct {
	ConsoleOptions
	SensitiveFields []string
}

func NewLog(opts Options) (LogClient, error) {
	logger, err := NewConsole(opts.ConsoleOptions)
	if err != nil {
		return nil, err
	}

	return &commonLogger{
		Console:         logger,
		sensitiveFields: opts.SensitiveFields,
	}, nil
}

type commonLogger struct {
	Console         *logrus.Logger
	sensitiveFields []string
}

type LogClient interface {
	Info(msg string, fields ...interface{})
	Trace(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	ErrorMessage(msg string, fields ...interface{})
	Error(err error, fields ...interface{})
	Fatal(err error, fields ...interface{})
	Panic(err error, fields ...interface{})
	Debug(msg string, fields ...interface{})

	SetOptionsToCtx(ctx context.Context, options ...options.LoggerOption) context.Context
	OptionsFromCtx(ctx context.Context) *options.LoggerOptions

	InfoCtx(ctx context.Context, msg string, fields ...interface{})
	TraceCtx(ctx context.Context, msg string, fields ...interface{})
	WarnCtx(ctx context.Context, msg string, fields ...interface{})
	ErrorMessageCtx(ctx context.Context, msg string, fields ...interface{})
	ErrorCtx(ctx context.Context, err error, fields ...interface{})
	FatalCtx(ctx context.Context, err error, fields ...interface{})
	PanicCtx(ctx context.Context, err error, fields ...interface{})
	DebugCtx(ctx context.Context, msg string, fields ...interface{})
}

func (l *commonLogger) Info(msg string, fields ...interface{}) {
	l.setFields(options.WithExtras(fields)).Info(msg)
}

func (l *commonLogger) Trace(msg string, fields ...interface{}) {
	l.setFields(options.WithExtras(fields)).Trace(msg)
}

func (l *commonLogger) Warn(msg string, fields ...interface{}) {
	l.setFields(options.WithExtras(fields)).Warn(msg)
}

func (l *commonLogger) Debug(msg string, fields ...interface{}) {
	l.setFields(options.WithExtras(fields)).Debug(msg)
}

func (l *commonLogger) ErrorMessage(msg string, fields ...interface{}) {
	l.setFields(options.WithException(&options.ExceptionData{
		Message: msg,
	})).Error()
}

func (l *commonLogger) Error(err error, fields ...interface{}) {
	l.setFields(options.WithException(&options.ExceptionData{
		Type:       fmt.Sprintf("%T", err),
		Message:    err.Error(),
		Stacktrace: l.formatLogStack(string(debug.Stack())),
	})).Error()
}

func (l *commonLogger) Fatal(err error, fields ...interface{}) {
	l.setFields(options.WithException(&options.ExceptionData{
		Type:       fmt.Sprintf("%T", err),
		Message:    err.Error(),
		Stacktrace: l.formatLogStack(string(debug.Stack())),
	})).Fatal()
}

func (l *commonLogger) Panic(err error, fields ...interface{}) {
	l.setFields(options.WithException(&options.ExceptionData{
		Type:       fmt.Sprintf("%T", err),
		Message:    err.Error(),
		Stacktrace: l.formatLogStack(string(debug.Stack())),
	})).Panic()
}

func (l *commonLogger) formatConsoleExtras(extras interface{}) any {
	if len(l.sensitiveFields) > 0 {
		formatConsoleExtrasWithObscure(l.sensitiveFields, extras)
	}
	return extras
}

func (l *commonLogger) setFields(optValues ...options.LoggerOption) *logrus.Entry {
	opts := &options.LoggerOptions{}
	for _, setOption := range optValues {
		setOption(opts)
	}
	if opts.Extras != nil {
		opts.Extras = l.formatConsoleExtras(opts.Extras)
	}
	return l.Console.
		WithField(LogOptionsField, *opts)
}
