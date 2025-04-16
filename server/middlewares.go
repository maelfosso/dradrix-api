package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
	"stockinos.com/api/services"
	"stockinos.com/api/storage"
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

func (s *Server) convertJWTTokenToMember(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		claims := ctx.Value(services.JwtClaimsKey)
		log.Println("convert jwt token to member : ", claims)
		if claims == nil {
			ctx = context.WithValue(ctx, services.JwtUserKey, nil)
			next.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		data := claims.(map[string]interface{})
		phoneNumber := data["PhoneNumber"].(string)

		user, err := s.database.Storage.GetUserByPhoneNumber(ctx, storage.GetUserByPhoneNumberParams{
			PhoneNumber: phoneNumber,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		s.log.Info(
			"Current Member",
			zap.Any("JWT user", user),
		)

		ctx = context.WithValue(ctx, services.JwtUserKey, user)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
