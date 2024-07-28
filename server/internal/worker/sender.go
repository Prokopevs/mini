package worker

import (
	"context"
	"time"

	"github.com/Prokopevs/mini/server/internal/service"
	"go.uber.org/zap"
)

// Producer - kafka producer.
type Messaging interface {
	NotifyMessagesCreated(context.Context, []int) error
	Close() error
}

type Sender struct {
	timeout       time.Duration
	db            service.DB
	logger        *zap.SugaredLogger
	kafkaProducer Messaging
	limit         int
}

func NewSender(timeout time.Duration, db service.DB, logger *zap.SugaredLogger, kafkaProducer Messaging, limit int) *Sender {
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
		select {
		case <-ctx.Done():
			t.logger.Info("Stopping process messages")
			ticker.Stop()
			t.kafkaProducer.Close()
			return
		case <-ticker.C:
		}

		t.ProcessMessages(ctx)
	}
}

func (t *Sender) ProcessMessages(ctx context.Context) {
	err := t.db.RunTx(ctx, func(ctx context.Context) error {
		// get messages with status='idle'
		ids, err := t.db.GetMessagesEvent(ctx, t.limit)
		if err != nil {
			t.logger.Errorw("failed to get messages", "err", err)
			return err
		}

		if len(ids) == 0 {
			return nil
		}

		// send to kafka
		err = t.kafkaProducer.NotifyMessagesCreated(ctx, ids)
		if err != nil {
			t.logger.Errorw("write to kafka", "err", err)
			return err
		}
		t.logger.Infow("Successfully sended messages to kafka", "ids", ids)

		// mark messages as sended
		status := "send"
		err = t.db.UpdateMessages(ctx, status, ids)
		if err != nil {
			t.logger.Errorw("failed update messages", "err", err)
			return err
		}
		t.logger.Infow("Changed statuses to `send` in db for", "ids", ids)

		return nil
	})
	if err != nil {
		t.logger.Errorw("tx processMessages error", "err", err)
	}
}
