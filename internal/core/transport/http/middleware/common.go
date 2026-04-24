package core_middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIdHeader = "X-Request-Id"

func RequestId() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIdHeader)
			if requestId == "" {
				requestId = uuid.NewString()
			}
			r.Header.Set(requestIdHeader, requestId)
			w.Header().Set(requestIdHeader, requestId)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIdHeader)

			logger := log.With(
				zap.String("request_id", requestId),
				zap.String("url", r.URL.String()),
				zap.String("method", r.Method),
			)

			ctx := context.WithValue(r.Context(), "logger", logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}

func PanicRecovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			responseHandler := response.NewHandlerResponse(log, w)
			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"during panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">> incoming request",
				zap.Time("time", time.Now().UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<< outgoing response",
				zap.Int("status code:", rw.GetStatus()),
				zap.Duration("latency", time.Now().Sub(before)),
			)
		})
	}
}
