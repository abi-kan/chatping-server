package core_http_middleware

import (
	core_logger "chatping-server/internal/core/logger"
	core_http_response "chatping-server/internal/core/transport/http/response"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
)

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]struct{}{
				"http://localhost:5051": {},
			}

			origin := r.Header.Get("Origin")
			if _, ok := allowedOrigins[origin]; ok {
				allowedMethods := strings.Join([]string{
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodPatch,
					http.MethodDelete,
					http.MethodOptions,
				}, ", ")
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			newLog := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)
			ctx := core_logger.ToContext(r.Context(), newLog)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := core_logger.FromContext(ctx)
			customRW := core_http_response.NewResponseWriter(w)

			startTime := time.Now().UTC()
			logger.Debug(
				">>> Incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", startTime),
			)
			next.ServeHTTP(customRW, r)
			logger.Debug(
				">>> Done HTTP request",
				zap.Int("status_code", customRW.GetStatusCode()),
				zap.Duration("latency", time.Since(startTime)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewResponseHandler(
				logger,
				w,
			)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"During hangling a HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
