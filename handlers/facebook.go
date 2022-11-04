package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func FacebookVerificationWebhooks(mux chi.Router) {
	mux.Get("/webhooks", func(w http.ResponseWriter, r *http.Request) {})
}
