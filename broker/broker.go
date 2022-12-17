package broker

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Broker struct {
	Nats    *nats.Conn
	log     *zap.Logger
	servers string
}

type Options struct {
	Servers string
	Log     *zap.Logger
}

func NewBroker(opts Options) *Broker {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	return &Broker{
		servers: opts.Servers,
		log:     opts.Log,
	}
}

func (b *Broker) Connect() error {
	b.log.Info("Connecting to NATS JetStream", zap.String("nats servers", b.servers))

	var err error
	b.Nats, err = nats.Connect(b.servers)
	if err != nil {
		b.log.Fatal("Failed to connect to the NATS JetStream server")
		return err
	}

	return nil
}
