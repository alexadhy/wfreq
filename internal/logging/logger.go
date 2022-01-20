package logging

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger and is able to create new package-aware loggers
type Logger struct {
	*logrus.Logger
}

// New creates a new package-aware logger with formatting string
func New() *Logger {
	l := logrus.New()
	e := logrus.NewEntry(l).Logger
	return &Logger{
		e,
	}
}

type ctxKeyLogger int

// LoggerKey defines logger HTTP context key.
const LoggerKey ctxKeyLogger = -1

// SetLoggerMiddleware sets logger to context of HTTP requests.
func SetLoggerMiddleware(log logrus.FieldLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if reqID := middleware.GetReqID(ctx); reqID != "" && log != nil {
				ctx = context.WithValue(ctx, LoggerKey, log.WithField("RequestID", reqID))
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
