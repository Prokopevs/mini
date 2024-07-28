package kafka

import (
	"context"
	"time"
	"strconv"
	"github.com/segmentio/kafka-go"
)

// producer - producer implementation.
type producer struct {
	writer *kafka.Writer
}

// NewProducer - ConsumerImpl constructor.
func NewProducer(topic string, brokers []string, groupID string) *producer {
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

func (p *producer) NotifyMessagesCreated(ctx context.Context, ids []int) error {
	kafkaMsgs := make([]kafka.Message, 0, len(ids))

	// create kafka messages array
	for _, id := range ids {
		idStr := strconv.Itoa(id)
		kafkaMsgs = append(kafkaMsgs, kafka.Message{
			Key:   []byte(idStr),
			Value: []byte(idStr),
		})
	}
	
	err := p.writer.WriteMessages(ctx, kafkaMsgs...)
	if err != nil {
		return err
	}

	return nil
}

func (p *producer) Close() error {
	return p.writer.Close()
}
