package publishers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"stockinos.com/api/broker"
	"stockinos.com/api/models"
)

type WhatsappMessageReceivedPublisher struct {
	// broker.BasePublisher
	js nats.JetStream
}

func (p *WhatsappMessageReceivedPublisher) stream() string {
	return "whatsapp:message"
}

func (p *WhatsappMessageReceivedPublisher) subject() broker.Subject {
	return broker.WhatsappMessageReceived
}

func (p *WhatsappMessageReceivedPublisher) Publish(data models.WhatsappMessage) error {
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

func NewWhatsappMessageReceivedPublisher(b broker.Broker) *WhatsappMessageReceivedPublisher {
	p := &WhatsappMessageReceivedPublisher{}

	err := b.BeforePublishing(
		p.stream(),
		p.subject().String(),
	)
	if err != nil {
		fmt.Println("NewWhatsappMessageReceivedPublisher - ", err)
		log.Fatalln(err)
	}

	p.js = b.JetStream
	return p
}
