package requests

type WhatsappSendMessageResponse struct {
	MessagingProduct string            `json:"messaging_product",omitemtpy`
	Contacts         []WhatsappContact `json:"contacts",omitempty`
	Messages         []WhatsappMessage `json:"messages",omitempty`
}

type WhatsappContact struct {
	Input string `json:"input",omitempty`
	WaId  string `json:"wa_id",omitempty`
}

type WhatsappMessage struct {
	ID string `json:"id",omitempty`
}

type SignInResult struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Token       string `json:"token,omitempty"`
}
