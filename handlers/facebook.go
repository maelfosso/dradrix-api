package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/utils"
)

func FacebookVerificationWebhooks(mux chi.Router) {
	mux.Get("/webhook", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		hubMode := query.Get("hub.mode")
		hubVerifyToken := query.Get("hub.verify_token")
		hubChallenge := query.Get("hub.challenge")

		if hubMode == "subscribe" && hubVerifyToken == utils.GetDefault("FACEBOOK_TOKEN", "stockinos-token") {
			w.Write([]byte(hubChallenge))
		}
	})

	mux.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("\nIncoming webhook: ", r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}
