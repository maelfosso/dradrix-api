package broker

import (
	"reflect"
)

type Event struct {
	subject Subject
	stream  string
	data    reflect.Type // func(reflect.Value) reflect.Value
}

// type WhatsappMessageReceivedEvent struct {
// 	subject Subject.
// 	data    models.WhatsappMessage
// }

// func NewWhatsappMessageReceivedEvent() *Event {
// 	return &Event{
// 		subject: WhatsappMessageReceived,
// 		stream:  "whatsapp:message",
// 		data:    models.WhatsappMessage,
// 	}
// }

// type WhatsappMessageReceivedEvent = NewWhatsappMessageReceivedEvent()
