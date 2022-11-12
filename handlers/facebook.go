package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/models"
	"stockinos.com/api/utils"
)

func FacebookWebhook(mux chi.Router) {
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
		body, _ := io.ReadAll(r.Body)
		fmt.Println("\nIncoming webhook: ", string(body))

		// decoder := json.NewDecoder(r.Body)
		buffer := bytes.NewBuffer(body)
		decoder := json.NewDecoder(buffer)
		var data models.WebhookData
		if err := decoder.Decode(&data); err != nil {
			fmt.Println("\nIncoming decode: ", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			return
		}

		fmt.Println("\nIncoming decoded: ", data)
		fmt.Println()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}
