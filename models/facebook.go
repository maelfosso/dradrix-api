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
	Value WebhookWhatsapp `json:"value",omitempty`
}

// WebhookEntryChangeValue Contains information about Whatsapp
// messaging_product = whatsapp
type WebhookWhatsapp struct { // Cont
	MessagingProduct string            `json:"messaging_product",omitemtpy`
	Metadata         WhatsappMetadata  `json:"metadata",omitemtpy`
	Contacts         []WhatsappContact `json:"contacts",omitempty`
	Messages         []WhatsappMessage `json:"messages",omitempty`
	Statuses         []WhatsappStatus  `json:"statuses",omitempty`
}

// Whatsapp Business account information
type WhatsappMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number",omitempty`
	PhoneNumberId      string `json:"phone_number_id",omitempty`
}

// Whatsapp user. For our case the client.
// Normally, we have to create an account for this user in our platform
type WhatsappContact struct {
	Profile WhatsappContactProfile `json:"profile",omitempty`
	WaId    string                 `json:"wa_id",omitempty`
}

type WhatsappContactProfile struct {
	Name string `json:"name",omitempty`
}

// The Whatsapp message sent by the client
type WhatsappMessage struct {
	gorm.Model
	ID                  string `json:"id",omitemtpy gorm:"primaryKey"`
	From                string `json:"from",omitempty`
	To                  string `json:"to",omitempty`
	Timestamp           string `json:"timestamp",omitempty`
	Type                string `json:"type",omitempty` // text, image, audio
	WhatsappMessageType `json:",inline" gorm:"embedded`
}

type WhatsappMessageType struct {
	TextID  uuid.UUID            `json:"text_id,omitempty" gorm:"default:null"`
	Text    WhatsappMessageText  `json:"text",omitempty`
	ImageID string               `json:"image_id,omitempty" gorm:"default:null"`
	Image   WhatsappMessageImage `json:"image",omitempty`
	AudioID string               `json:"audio_id,omitempty" gorm:"default:null"`
	Audio   WhatsappMessageAudio `json:"audio",omitempty`
}

type WhatsappMessageText struct {
	gorm.Model
	ID   uuid.UUID `json:"id,omitempty" gorm:"type:uuid;primary_key;"`
	Body string    `json:"body",omitempty`
}

type WhatsappMessageImage struct {
	gorm.Model
	Caption  string `json:"caption",omitemtpy`
	MimeType string `json:"mime_type",omitempty` // image/jpeg,
	Sha256   string `json:"sha256",omitemtpy`
	ID       string `json:"id",omitempty gorm:"primaryKey"`
}

type WhatsappMessageAudio struct {
	gorm.Model
	MimeType string `json:"mime_type",omitempty` // audio/ogg; codecs=opus -
	Sha256   string `json:"sha256",omitemtpy`
	ID       string `json:"id",omitempty gorm:"primaryKey"`
	Voice    bool   `json:"voice",omitempty`
}

type WhatsappStatus struct {
	gorm.Model
	ID           string                     `json:"id",omitempty gorm:"primaryKey"`
	Status       string                     `json:"status",omitempty`
	Timestamp    string                     `json:"timestamp",omitempty`
	RecipientId  string                     `json:"recipient_id",omitempty`
	Conversation WhatsappStatusConversation `json:"conversation",omitempty`
	Pricing      WhatsappStatusPricing      `json:"pricing",omitempty`
}

type WhatsappStatusConversation struct {
	ID                  string                           `json:"id",omitempty gorm:"primaryKey"`
	ExpirationTimestamp string                           `json:"expiration_timestamp",omitempty`
	Origin              WhatsappStatusConversationOrigin `json:"origin",omitempty`
}

type WhatsappStatusConversationOrigin struct {
	Type string `json:"type",omitempty`
}

type WhatsappStatusPricing struct {
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
