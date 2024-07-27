package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer - kafka producer.
type Producer interface {
	SendMessages(context.Context, []kafka.Message) error
	Close() error
}

// producer - producer implementation.
type producer struct {
	writer *kafka.Writer
}

// NewProducer - ConsumerImpl constructor.
func NewProducer(topic string, brokers []string, groupID string) Producer {
	return &producer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
			Dialer: &kafka.Dialer{
				Timeout:   10 * time.Second,
				DualStack: true,
			},
		}),
	}
}

func (p *producer) SendMessages(ctx context.Context, msg []kafka.Message) error {
	err := p.writer.WriteMessages(ctx, msg...)
	if err != nil {
		return err
	}

	return nil
}

func (p *producer) Close() error {
	return p.writer.Close()
}
