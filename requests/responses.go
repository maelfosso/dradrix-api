package requests

type WhatsAppSendTextMessageResponse struct {
	MessagingProduct string            `json:"messaging_product",omitemtpy`
	Contacts         []WhatsAppContact `json:"contacts",omitempty`
	Messages         []WhatsAppMessage `json:"messages",omitempty`
}

type WhatsAppContact struct {
	Input string `json:"input",omitempty`
	WaId  string `json:"wa_id",omitempty`
}

type WhatsAppMessage struct {
	ID string `json:"id",omitempty`
}
