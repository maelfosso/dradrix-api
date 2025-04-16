package broker

type IBaseSubscriber interface {
	Subscribe(data interface{}) error
}

type BaseSubscriber struct {
	stream   func() string
	subject  func() Subject
	parseMsg func()
	broker   Broker
}

func NewSubscriber(broker Broker) *BaseSubscriber {
	return &BaseSubscriber{
		broker: broker,
	}
}
