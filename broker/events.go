package broker

import (
	"reflect"
)

type Event struct {
	subject Subject
	stream  string
	data    reflect.Type // func(reflect.Value) reflect.Value
}

// type WhatsAppMessageReceivedEvent struct {
// 	subject Subject.
// 	data    models.WhatsAppMessage
// }

// func NewWhatsAppMessageReceivedEvent() *Event {
// 	return &Event{
// 		subject: WhatsAppMessageReceived,
// 		stream:  "whatsapp:message",
// 		data:    models.WhatsAppMessage,
// 	}
// }

// type WhatsAppMessageReceivedEvent = NewWhatsAppMessageReceivedEvent()
