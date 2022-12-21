package publishers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"stockinos.com/api/broker"
	"stockinos.com/api/models"
)

type WhatsAppMessageReceivedPublisher struct {
	// broker.BasePublisher
	js nats.JetStream
}

func (p *WhatsAppMessageReceivedPublisher) stream() string {
	return "whatsapp:message"
}

func (p *WhatsAppMessageReceivedPublisher) subject() broker.Subject {
	return broker.WhatsAppMessageReceived
}

func (p *WhatsAppMessageReceivedPublisher) Publish(data models.WhatsAppMessage) error {
	fmt.Println("\nMessage to Publish: ", data)
	fmt.Println()
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}

	pubAck, err := p.js.Publish(p.subject().String(), b)
	if err != nil {
		return fmt.Errorf("error publishing message: %w", err)
	}
	fmt.Println("Pub Ack : ", pubAck)
	fmt.Println()

	return nil
}

func NewWhatsAppMessageReceivedPublisher(b broker.Broker) *WhatsAppMessageReceivedPublisher {
	p := &WhatsAppMessageReceivedPublisher{}

	err := b.BeforePublishing(
		p.stream(),
		p.subject().String(),
	)
	if err != nil {
		fmt.Println("NewWhatsAppMessageReceivedPublisher - ", err)
		log.Fatalln(err)
	}

	p.js = b.JetStream
	return p
}
