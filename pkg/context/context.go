package marketcontext

import (
	"context"
	"net/http"

	"github.com/osalomon89/test-crud-api/pkg/log"
)

const RequestIDKey = "X-Request-Id"

type loggerKey struct{}

func New(request *http.Request) context.Context {
	requestID := request.Header.Get(RequestIDKey)
	if len(requestID) == 0 {
		return context.WithValue(request.Context(), loggerKey{}, log.DefaultLogger())
	}

	return context.WithValue(request.Context(), loggerKey{}, log.NewLogger(requestID))
}

func Logger(ctx context.Context) log.ILogger {
	logger, ok := ctx.Value(loggerKey{}).(log.ILogger)
	if !ok {
		return log.DefaultLogger()
	}

	return logger
}
