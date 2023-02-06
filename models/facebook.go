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
	Statuses         []WhatsAppStatus  `json:"statuses",omitempty`
}

// Whatsapp Business account information
type WhatsAppMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number",omitempty`
	PhoneNumberId      string `json:"phone_number_id",omitempty`
}

// Whatsapp user. For our case the client.
// Normally, we have to create an account for this user in our platform
type WhatsAppContact struct {
	Profile WhatsAppContactProfile `json:"profile",omitempty`
	WaId    string                 `json:"wa_id",omitempty`
}

type WhatsAppContactProfile struct {
	Name string `json:"name",omitempty`
}

// The Whatsapp message sent by the client
type WhatsAppMessage struct {
	gorm.Model
	ID                  string `json:"id",omitemtpy gorm:"primaryKey"`
	From                string `json:"from",omitempty`
	To                  string `json:"to",omitempty`
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

type WhatsAppStatus struct {
	gorm.Model
	ID           string                     `json:"id",omitempty gorm:"primaryKey"`
	Status       string                     `json:"status",omitempty`
	Timestamp    string                     `json:"timestamp",omitempty`
	RecipientId  string                     `json:"recipient_id",omitempty`
	Conversation WhatsAppStatusConversation `json:"conversation",omitempty`
	Pricing      WhatsAppStatusPricing      `json:"pricing",omitempty`
}

type WhatsAppStatusConversation struct {
	ID                  string                           `json:"id",omitempty gorm:"primaryKey"`
	ExpirationTimestamp string                           `json:"expiration_timestamp",omitempty`
	Origin              WhatsAppStatusConversationOrigin `json:"origin",omitempty`
}

type WhatsAppStatusConversationOrigin struct {
	Type string `json:"type",omitempty`
}

type WhatsAppStatusPricing struct {
	Billable     bool   `json:"billable",omitempty`
	PricingModel string `json:"pricing_model",omitempty`
	Category     string `json:"category",omitempty`
}

// the message has been sent
// {"object":"whatsapp_business_account","entry":[{"id":"101270969527843","changes":[{"value":{"messaging_product":"whatsapp","metadata":{"display_phone_number":"237620388204","phone_number_id":"115426628094550"},"statuses":[{"id":"wamid.HBgMMjM3Njk1MTY1MDMzFQIAERgSOUFFOTEyMDk1OEFGQzgwMTMxAA==","status":"sent","timestamp":"1672720505","recipient_id":"237695165033","conversation":{"id":"687c7adf6abb01537e9c1d1bf2cdf3e7","expiration_timestamp":"1672806960","origin":{"type":"business_initiated"}},"pricing":{"billable":true,"pricing_model":"CBP","category":"business_initiated"}}]},"field":"messages"}]}]}

// the message has been delivered
// {"object":"whatsapp_business_account","entry":[{"id":"101270969527843","changes":[{"value":{"messaging_product":"whatsapp","metadata":{"display_phone_number":"237620388204","phone_number_id":"115426628094550"},"statuses":[{"id":"wamid.HBgMMjM3Njk1MTY1MDMzFQIAERgSOUFFOTEyMDk1OEFGQzgwMTMxAA==","status":"delivered","timestamp":"1672720507","recipient_id":"237695165033","conversation":{"id":"687c7adf6abb01537e9c1d1bf2cdf3e7","origin":{"type":"business_initiated"}},"pricing":{"billable":true,"pricing_model":"CBP","category":"business_initiated"}}]},"field":"messages"}]}]}

// the message has been read
// {"object":"whatsapp_business_account","entry":[{"id":"101270969527843","changes":[{"value":{"messaging_product":"whatsapp","metadata":{"display_phone_number":"237620388204","phone_number_id":"115426628094550"},"statuses":[{"id":"wamid.HBgMMjM3Njk1MTY1MDMzFQIAERgSOUFFOTEyMDk1OEFGQzgwMTMxAA==","status":"read","timestamp":"1672720782","recipient_id":"237695165033"}]},"field":"messages"}]}]}

// a message has different state at different times

// The user has read the message sent from Whatsapp Messager API Platform
// {
//   "object": "whatsapp_business_account",
//   "entry": [
//     {
//       "id": "101270969527843",
//       "changes": [
//         {
//           "value": {
//             "messaging_product": "whatsapp",
//             "metadata": {
//               "display_phone_number": "237620388204",
//               "phone_number_id": "115426628094550"
//             },
//             "statuses": [
//               {
//                 "id": "wamid.HBgMMjM3Njc4OTA4OTg5FQIAERgSNzBBOTFFMTM4MjEzMzFBQjE0AA==",
//                 "status": "read",
//                 "timestamp": "1672717775",
//                 "recipient_id": "237678908989"
//               }
//             ]
//           },
//           "field": "messages"
//         }
//       ]
//     }
//   ]
// }
