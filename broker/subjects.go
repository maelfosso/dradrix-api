package broker

type Subject string

const (
	WhatsAppMessageReceived Subject = "whatsapp:message:received"
	MessageWoZSent          Subject = "whatsapp:message:woz:sent"
)

func (s Subject) String() string {
	switch s {
	case WhatsAppMessageReceived:
		return "whatsapp:message:received"
	case MessageWoZSent:
		return "whatsapp:message:woz:sent"
	default:
		return "unknown"
	}
}
