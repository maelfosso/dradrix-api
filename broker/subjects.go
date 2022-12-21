package broker

type Subject string

const (
	WhatsAppMessageReceived Subject = "whatsapp:message:received"
)

func (s Subject) String() string {
	switch s {
	case WhatsAppMessageReceived:
		return "whatsapp:message:received"
	default:
		return "unknown"
	}
}
