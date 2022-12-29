package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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
	gorm.Model
	ID                  string `json:"id",omitemtpy gorm:"primaryKey"`
	From                string `json:"from",omitempty`
	Timestamp           string `json:"timestamp",omitempty`
	Type                string `json:"type",omitempty` // text, image, audio
	WhatsAppMessageType `json:",inline" gorm:"embedded`
}

type WhatsAppMessageType struct {
	TextID  uuid.UUID            `json:"text_id,omitempty" gorm:"default:null"`
	Text    WhatsAppMessageText  `json:"text",omitempty`
	ImageID string               `json:"image_id,omitempty" gorm:"default:null"`
	Image   WhatsAppMessageImage `json:"image",omitempty`
	AudioID string               `json:"audio_id,omitempty" gorm:"default:null"`
	Audio   WhatsAppMessageAudio `json:"audio",omitempty`
}

type WhatsAppMessageText struct {
	gorm.Model
	ID   uuid.UUID `json:"id,omitempty" gorm:"type:uuid;primary_key;"`
	Body string    `json:"body",omitempty`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (w *WhatsAppMessageText) BeforeCreate(tx *gorm.DB) error {
	w.ID = uuid.NewV4()
	return nil
}

type WhatsAppMessageImage struct {
	gorm.Model
	Caption  string `json:"caption",omitemtpy`
	MimeType string `json:"mime_type",omitempty` // image/jpeg,
	Sha256   string `json:"sha256",omitemtpy`
	ID       string `json:"id",omitempty gorm:"primaryKey"`
}

type WhatsAppMessageAudio struct {
	gorm.Model
	MimeType string `json:"mime_type",omitempty` // audio/ogg; codecs=opus -
	Sha256   string `json:"sha256",omitemtpy`
	ID       string `json:"id",omitempty gorm:"primaryKey"`
	Voice    bool   `json:"voice",omitempty`
}
