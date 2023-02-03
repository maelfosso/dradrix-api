package broker

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Broker struct {
	Nats      *nats.Conn
	JetStream nats.JetStreamContext
	log       *zap.Logger
	servers   string
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

func (b *Broker) Setup() error {
	b.log.Info("Setting up NATS JetStream")

	var err error
	b.JetStream, err = b.Nats.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		b.log.Fatal("Failed to create jet stream: ", zap.Error(err))
	}

	return nil
}

// // createStream creates a stream by using JetStreamContext
// func createStream(js nats.JetStreamContext) error {
// 	// Check if the ORDERS stream already exists; if not, create it.
// 	stream, err := js.StreamInfo(streamName)
// 	if err != nil {
// 		 log.Println(err)
// 	}
// 	if stream == nil {
// 		 log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
// 		 _, err = js.AddStream(&nats.StreamConfig{
// 				Name:     streamName,
// 				Subjects: []string{streamSubjects},
// 		 })
// 		 if err != nil {
// 				return err
// 		 }
// 	}
// 	return nil
// }

func (b *Broker) CheckStreamOrCreate(name string) (*nats.StreamInfo, error) {
	var stream *nats.StreamInfo
	stream, err := b.JetStream.StreamInfo(name)
	if err != nil {
		// return nil, fmt.Errorf("error getting stream [%s] info: %w", name, err)
		log.Println(err)
	}
	if stream == nil {
		stream, err = b.JetStream.AddStream(&nats.StreamConfig{
			Name:     name,
			Storage:  nats.FileStorage,
			Subjects: []string{fmt.Sprint(stream, ":*")}, // stream:*
		})

		if err != nil {
			return nil, fmt.Errorf("error creating stream [%s] info: %w", name, err)
		}
	}

	return stream, nil
}

func (b *Broker) AddSubjectToStream(stream *nats.StreamInfo, subject string) (*nats.StreamInfo, error) {
	log.Println("AddSugjectToStream - Stream ", stream, stream.Config)
	var err error

	// Check if the subject is already into the stream
	for _, s := range stream.Config.Subjects {
		if s == subject {
			return stream, nil
		}
	}

	stream, err = b.JetStream.UpdateStream(&nats.StreamConfig{
		Name:     stream.Config.Name,
		Subjects: append(stream.Config.Subjects, subject),
	})
	log.Println("After UPdateStream ", stream, err)
	if err != nil {
		return nil, fmt.Errorf("error adding [%s] to the stream [%s] info: %w", subject, stream.Config.Name, err)
	}

	return stream, err
}

func (b *Broker) BeforePublishing(streamName, subjectName string) error {
	stream, err := b.CheckStreamOrCreate(streamName)
	if err != nil {
		return fmt.Errorf("error when checking or creating [%s]: %w", streamName, err)
	}
	fmt.Println("Before Publishing : Stream created - ", stream, err)

	subjectAlreadyPresent := false
	for _, v := range stream.Config.Subjects {
		if v == streamName {
			subjectAlreadyPresent = true

			break
		}
	}

	if !subjectAlreadyPresent {
		_, err = b.AddSubjectToStream(stream, subjectName)
		if err != nil {
			return fmt.Errorf("error when adding subject [%s] to stream [%s]: %w", subjectName, streamName, err)
		}
	}

	return nil
}
