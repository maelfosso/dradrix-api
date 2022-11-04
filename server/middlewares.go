package server

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewStatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (srw *statusResponseWriter) WriteHeader(statusCode int) {
	srw.statusCode = statusCode
	srw.ResponseWriter.WriteHeader(statusCode)
}

func (s *Server) requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		srw := NewStatusResponseWriter(w)

		defer func() {
			s.log.Info(
				"Request sent",
				zap.String("method", r.Method),
				zap.Duration("started at", time.Since(start)),
				zap.Int("status", srw.statusCode),
				zap.String("host", r.Host),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
			)
		}()

		next.ServeHTTP(srw, r)
	})
}
