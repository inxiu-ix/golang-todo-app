package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/inxiu-ix/golang-todo-app/internal/core/logger"
	core_http_response "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
)

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]struct{}{
				"http://localhost:5050": {},
			}

			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == "OPTIONS" {
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

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			before := time.Now().UTC()
			rw := core_http_response.NewResponseWriter(w)

			log.Debug(
				">>>>>> incoming HTTP request",
				zap.String("method", r.Method),
				zap.Time("time", before.UTC()))

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("duration", time.Now().UTC().Sub(before)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				ctx := r.Context()
				log := core_logger.FromContext(ctx)
				responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during request processing occurred unexpected panic")
				}

				next.ServeHTTP(w, r)
			}()
		})
	}
}
