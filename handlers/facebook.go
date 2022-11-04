package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func VerificationWebhooks(mux chi.Router) {
	mux.Get("/webhooks", func(w http.ResponseWriter, r *http.Request) {})
}
