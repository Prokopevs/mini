package worker

import (
	"context"
	"strconv"
	"time"

	kf "github.com/Prokopevs/mini/server/internal/kafka"
	"github.com/Prokopevs/mini/server/internal/service"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Sender struct {
	timeout       time.Duration
	db            service.DB
	logger        *zap.SugaredLogger
	kafkaProducer kf.Producer
	limit         int
}

func NewSender(timeout time.Duration, db service.DB, logger *zap.SugaredLogger, kafkaProducer kf.Producer, limit int) *Sender {
	return &Sender{
		timeout:       timeout,
		db:            db,
		logger:        logger,
		kafkaProducer: kafkaProducer,
		limit:         limit,
	}
}

func (t *Sender) RunProcessingMessages(ctx context.Context) {
	ticker := time.NewTicker(t.timeout)

	t.logger.Info("Run process messages")

	for {
		kafkaMsgs := make([]kafka.Message, 0, t.limit)

		select {
		case <-ctx.Done():
			t.logger.Info("Stopping process messages")
			ticker.Stop()
			t.kafkaProducer.Close()
			return
		case <-ticker.C:
		}

		// get messages with status='idle'
		ids, err := t.db.GetMessagesEvent(ctx, t.limit)
		if err != nil {
			t.logger.Infow("failed to get messages", "err", err)
			continue
		}

		if len(ids) == 0 {
			continue
		}

		// create kafka messages array
		for _, id := range ids {
			idStr := strconv.Itoa(id)
			kafkaMsgs = append(kafkaMsgs, kafka.Message{
				Value: []byte(idStr),
			})
		}

		// send to kafka
		err = t.kafkaProducer.SendMessages(ctx, kafkaMsgs)
		if err != nil {
			t.logger.Infow("write to kafka", "err", err)
			continue
		}

		// mark messages as sended
		status := "send"
		err = t.db.UpdateMessages(ctx, status, ids)
		if err != nil {
			t.logger.Infow("failed update messages", "err", err)
			continue
		}
	}
}
