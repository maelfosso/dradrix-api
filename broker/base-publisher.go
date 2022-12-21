package broker

type IBasePublisher interface {
	Publish(data interface{}) error
}

type BasePublisher struct {
	stream  func() string
	subject func() Subject
	broker  Broker
	// Js      nats.JetStream
}

func NewPublisher(broker Broker) *BasePublisher {
	return &BasePublisher{
		broker: broker,
		// js:     broker.JetStream,
	}
}

func (b *BasePublisher) Init() error {
	stream, err := b.broker.CheckStreamOrCreate(b.stream())
	if err != nil {
		return err
	}

	_, err = b.broker.AddSubjectToStream(stream, b.subject().String())
	if err != nil {
		return err
	}

	return nil
}
