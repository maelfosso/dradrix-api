package subscribers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"stockinos.com/api/broker"
)

type WhatsAppMessageWoZSentData struct {
	Id      string `json:"id",omitemtpy`
	From    string `json:"from",omitemtpy`
	To      string `json:"to",omitemtpy`
	Message string `json:"message",omitemtpy`
}

type WhatsAppMessageWoZSentSubscriber struct {
	js nats.JetStream
}

func (p *WhatsAppMessageWoZSentSubscriber) stream() string {
	return "whatsapp:message"
}

func (p *WhatsAppMessageWoZSentSubscriber) subject() broker.Subject {
	return broker.WhatsAppMessageWoZSent
}

func (p *WhatsAppMessageWoZSentSubscriber) parseMsg(msg *nats.Msg) (*WhatsAppMessageWoZSentData, error) {
	log.Println("WhatsAppMessageWoZSentData : ", msg.Data, string(msg.Data[:]))

	var data WhatsAppMessageWoZSentData
	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshilling WhatsAppMessageWoZSentData : %s", err)
	}
	log.Println("WhatsAppMessageWoZSentData Parsed: ", data)

	return &data, nil
}

func (p *WhatsAppMessageWoZSentSubscriber) Subscribe() error {
	p.js.Subscribe(p.subject().String(), func(msg *nats.Msg) {
		data, err := p.parseMsg(msg)
		if err != nil {
			log.Println(err)
			// msg.Nak()
			msg.Ack()
		} else {
			log.Printf("monitor service subscribes from subject:%s\n", msg.Subject)
			log.Printf("To:%s, From: %s, Message:%s\n", data.To, data.From, data.Message)
		}

	}, nats.Durable("monitor"), nats.ManualAck())

	return nil
}

func NewWhatsAppMessageWoZSentSubscriber(b broker.Broker) *WhatsAppMessageWoZSentSubscriber {
	s := &WhatsAppMessageWoZSentSubscriber{}

	err := b.BeforePublishing(
		s.stream(),
		s.subject().String(),
	)
	if err != nil {
		fmt.Println("NewWhatsAppMessageWoZSentSubscriber - ", err)
		log.Fatalln(err)
	}

	s.js = b.JetStream
	return s
}
