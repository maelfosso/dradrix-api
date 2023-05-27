package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"stockinos.com/api/broker/publishers"
	"stockinos.com/api/handlers"
	"stockinos.com/api/storage"
)

type facebookWebhookStruct struct {
	*storage.Database
	*publishers.WhatsappMessageReceivedPublisher
}

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

	handlers.Root(s.mux)
	handlers.Health(s.mux)

	handlers.FacebookWebhook(s.mux, facebookWebhookStruct{
		Database:                         s.database,
		WhatsappMessageReceivedPublisher: publishers.NewWhatsappMessageReceivedPublisher(*s.nats),
	})

	s.mux.Group(func(r chi.Router) {

		// Auth
		r.Route("/auth", func(r chi.Router) {
			handlers.GetOTP(r, s.database)
			handlers.CheckOTP(r, s.database)
			handlers.ResendOTP(r, s.database)
		})

	})
}
