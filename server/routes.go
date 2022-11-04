package server

import (
	"github.com/go-chi/cors"
	"stockinos.com/api/handlers"
)

func (s *Server) setupRoutes() {
	s.mux.Use(s.requestLoggerMiddleware)
	s.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	handlers.Health(s.mux)
	handlers.FacebookVerificationWebhooks(s.mux)
}
