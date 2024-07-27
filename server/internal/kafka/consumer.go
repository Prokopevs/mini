package kafka

import (
	"context"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	
)

// Consumer - kafka consumer.
type Consumer interface {
	ReadMessage(context.Context) (CommitFunc, int, error)
	Close() error
}

// ConsumerImpl - consumer implementation.
type ConsumerImpl struct {
	reader *kafka.Reader
}

// NewComcumerImpl - ConsumerImpl constructor.
func NewConsumer(topic string, brokers []string, groupID string) *ConsumerImpl {
	return &ConsumerImpl{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
			Dialer: &kafka.Dialer{
				Timeout:   10 * time.Second,
				DualStack: true,
			},
		}),
	}
}

// UnmarshalMessage - fetch message from kafka broker.
func (c *ConsumerImpl) ReadMessage(ctx context.Context) (CommitFunc, int, error) {
	kafkaMsg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return func(ctx2 context.Context) error { return nil }, 0, err
	}

	id, err := strconv.Atoi(string(kafkaMsg.Value))
	if err != nil {
		return func(ctx2 context.Context) error { return nil }, 0, err
	}

	return func(ctx2 context.Context) error {
		return c.reader.CommitMessages(ctx2, kafkaMsg)
	}, id, nil
}

func (c *ConsumerImpl) Close() error {
	return c.reader.Close()
}
