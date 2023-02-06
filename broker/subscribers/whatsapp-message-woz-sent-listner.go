package subscribers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"stockinos.com/api/broker"
	"stockinos.com/api/models"
	"stockinos.com/api/services"
	"stockinos.com/api/storage"
)

type MessageWoZSentData struct {
	Id      string `json:"id",omitemtpy`
	From    string `json:"from",omitemtpy`
	To      string `json:"to",omitemtpy`
	Message string `json:"message",omitemtpy`
}

type MessageWoZSentSubscriber struct {
	js nats.JetStream
}

func (p *MessageWoZSentSubscriber) stream() string {
	return "whatsapp:message"
}

func (p *MessageWoZSentSubscriber) subject() broker.Subject {
	return broker.MessageWoZSent
}

func (p *MessageWoZSentSubscriber) parseMsg(msg *nats.Msg) (*models.WhatsAppMessage, error) {
	log.Println("MessageWoZSentData : ", msg.Data, string(msg.Data[:]))

	var data models.WhatsAppMessage // MessageWoZSentData
	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshilling MessageWoZSentData : %s", err)
	}
	log.Println("MessageWoZSentData Parsed: ", data)

	return &data, nil
}

func (p *MessageWoZSentSubscriber) Subscribe(db storage.Database) error {
	p.js.Subscribe(p.subject().String(), func(msg *nats.Msg) {
		data, err := p.parseMsg(msg)
		if err != nil {
			log.Println(err)
			// msg.Nak()
			msg.Ack()
		} else {
			log.Printf("monitor service subscribes from subject:%s\n", msg.Subject)
			log.Printf("To:%s, From: %s, Message:%s\n", data.To, data.From, data.Text.Body)

			err = services.HandleMessageSentByWoZ(db, *data)
			if err != nil {
				msg.Nak()
			} else {
				msg.Ack()
			}
		}

	}, nats.Durable("monitor"), nats.ManualAck())

	return nil
}

func NewMessageWoZSentSubscriber(b broker.Broker) *MessageWoZSentSubscriber {
	s := &MessageWoZSentSubscriber{}

	err := b.BeforePublishing(
		s.stream(),
		s.subject().String(),
	)
	if err != nil {
		fmt.Println("NewMessageWoZSentSubscriber - ", err)
		log.Fatalln(err)
	}

	s.js = b.JetStream
	return s
}
