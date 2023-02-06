package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/models"
	"stockinos.com/api/utils"
)

type facebookWebhookInterface interface {
	Publish(message models.WhatsAppMessage) error
	SaveWAMessages(ctx context.Context, messages []models.WhatsAppMessage) error
	SaveWAStatus(ctx context.Context, statuses []models.WhatsAppStatus) error
}

func FacebookWebhook(mux chi.Router, s facebookWebhookInterface) {
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

		// fmt.Println("\nIncoming decoded: ", data)
		// fmt.Println()

		changes := data.Entry[0].Changes[0].Value
		if len(changes.Messages) > 0 {
			messages := changes.Messages
			fmt.Println("\nMessage to send: ", messages[0])
			fmt.Println()
			err := s.Publish(messages[0])
			if err != nil {
				fmt.Println(err)
			}
			// s.SaveWAMessages(r.Context(), messages)
		}
		if len(changes.Statuses) > 0 {
			fmt.Println("\n**** STATUSES **** ")
			fmt.Println(changes.Statuses)
			fmt.Println()

			s.SaveWAStatus(r.Context(), changes.Statuses)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}
