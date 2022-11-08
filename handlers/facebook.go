package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/utils"
)

// WebhookBody
// object = whatsapp_business_account
type WebhookData struct {
	Object string         `json:"object",omitempty`
	Entry  []WebhookEntry `json:"entry",omitempty`
}

type WebhookEntry struct {
	ID      string               `json:"id",omitempty`
	Changes []WebhookEntryChange `json:"changes",omitempty`
}

// WebhookEntryChange
// field = messages for messages
type WebhookEntryChange struct {
	Field string          `json:"field",omitempty`
	Value WebhookWhatsApp `json:"value",omitempty`
}

// WebhookEntryChangeValue Contains information about WhatsApp
// messaging_product = whatsapp
type WebhookWhatsApp struct { // Cont
	MessagingProduct string            `json:"messaging_product",omitemtpy`
	Metadata         WhatsAppMetadata  `json:"metadata",omitemtpy`
	Contacts         []WhatsAppContact `json:"contacts",omitempty`
	Messages         []WhatsAppMessage `json:"messages",omitempty`
}

type WhatsAppMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number",omitempty`
	PhoneNumberId      string `json:"phone_number_id",omitempty`
}

type WhatsAppContact struct {
	Profile WhatsAppContactProfile `json:"profile",omitempty`
	WaId    string                 `json:"wa_id",omitempty`
}

type WhatsAppContactProfile struct {
	Name string `json:"name",omitempty`
}

type WhatsAppMessage struct {
	From                string `json:"from",omitempty`
	Id                  string `json:"id",omitemtpy`
	Timestamp           string `json:"timestamp",omitempty`
	Type                string `json:"type",omitempty` // text, image, audio
	WhatsAppMessageType `json:",inline"`
}

type WhatsAppMessageType struct {
	Text  WhatsAppMessageText  `json:"text",omitempty`
	Image WhatsAppMessageImage `json:"image",omitempty`
	Audio WhatsAppMessageAudio `json:"audio",omitempty`
}

type WhatsAppMessageText struct {
	Body string `json:"body",omitempty`
}

type WhatsAppMessageImage struct {
	Caption  string `json:"caption",omitemtpy`
	MimeType string `json:"mime_type",omitempty` // image/jpeg,
	Sha256   string `json:"sha256",omitemtpy`
	Id       string `json:"id",omitempty`
}

type WhatsAppMessageAudio struct {
	MimeType string `json:"mime_type",omitempty` // audio/ogg; codecs=opus -
	Sha256   string `json:"sha256",omitemtpy`
	Id       string `json:"id",omitempty`
	Voice    bool   `json:"voice",omitempty`
}

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
		var data WebhookData
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
