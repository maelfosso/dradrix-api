package broker

type Subject string

const (
	WhatsAppMessageReceived Subject = "whatsapp:message:received"
	WhatsAppMessageWoZSent  Subject = "whatsapp:message:woz:sent"
)

func (s Subject) String() string {
	switch s {
	case WhatsAppMessageReceived:
		return "whatsapp:message:received"
	case WhatsAppMessageWoZSent:
		return "whatsapp:message:woz:sent"
	default:
		return "unknown"
	}
}
