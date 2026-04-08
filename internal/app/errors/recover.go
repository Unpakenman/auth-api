package errors

import (
	"context"
	"fmt"
)

func RecoverHandler(ctx context.Context, logger LogClient) func() {
	return func() {
		if panicErr := recover(); panicErr != nil {
			err := fmt.Errorf("recovered from panic: %v", panicErr)
			if logger != nil {
				logger.ErrorCtx(ctx, err, panicErr)
			} else {
				fmt.Println(err, panicErr)
			}
		}
	}
}

type LogClient interface {
	ErrorCtx(ctx context.Context, err error, fields ...interface{})
}
