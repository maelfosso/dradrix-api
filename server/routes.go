package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"stockinos.com/api/broker/publishers"
	"stockinos.com/api/handlers"
	"stockinos.com/api/services"
	"stockinos.com/api/storage"
)

type facebookWebhookStruct struct {
	*storage.Database
	*publishers.WhatsappMessageReceivedPublisher
}

func (s *Server) setupRoutes() {
	appHandler := handlers.NewAppHandler()

	s.mux.Use(s.requestLoggerMiddleware)
	s.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.mux.Use(services.Verifier)
	s.mux.Use(services.ParseJwtToken)

	// Protected Routes
	s.mux.Group(func(r chi.Router) {
		r.Use(services.Authenticator)
		r.Use(s.convertJWTTokenToMember)

		handlers.GetCurrentUser(r)

		r.Route("/companies", func(r chi.Router) {
			appHandler.GetAllCompanies(r, s.database.Storage)
			appHandler.CreateCompany(r, s.database.Storage)

			r.Route("/{companyId}", func(r chi.Router) {
				appHandler.CompanyMiddleware(r, s.database.Storage)

				appHandler.GetCompany(r, s.database.Storage)
				appHandler.UpdateCompany(r, s.database.Storage)
				appHandler.DeleteCompany(r, s.database.Storage)

				r.Route("/activities", func(r chi.Router) {
					appHandler.GetAllActivities(r, s.database.Storage)
					appHandler.CreateActivity(r, s.database.Storage)

					r.Route("/{activityId}", func(r chi.Router) {
						appHandler.ActivityMiddleware(r, s.database.Storage)

						appHandler.GetActivity(r, s.database.Storage)
						appHandler.DeleteActivity(r, s.database.Storage)
						appHandler.UpdateCompany(r, s.database.Storage)

						r.Route("/data", func(r chi.Router) {
							appHandler.CreateData(r, s.database.Storage)
							appHandler.GetAllData(r, s.database.Storage)

							r.Route("{dataId}", func(r chi.Router) {
								appHandler.DataMiddleware(r, s.database.Storage)
							})
						})
					})
				})
			})
		})
	})

	// Public Routes
	s.mux.Group(func(r chi.Router) {
		handlers.Root(r)
		handlers.Health(r, s.database)

		r.Group(func(r chi.Router) {
			// Auth
			r.Route("/auth", func(r chi.Router) {
				handlers.CreateOTP(r, s.database.Storage)
				handlers.CheckOTP(r, s.database.Storage)
				// handlers.ResendOTP(r, s.database)
			})
		})
	})

	// handlers.FacebookWebhook(s.mux, facebookWebhookStruct{
	// 	Database:                         s.database,
	// 	WhatsappMessageReceivedPublisher: publishers.NewWhatsappMessageReceivedPublisher(*s.nats),
	// })
}
