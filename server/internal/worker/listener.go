package worker

import (
	"context"
	"fmt"

	kf "github.com/Prokopevs/mini/server/internal/kafka"
	"github.com/Prokopevs/mini/server/internal/service"
	"go.uber.org/zap"
)

type Listener struct {
	logger        *zap.SugaredLogger
	db            service.DB
	kafkaConsumer kf.Consumer
}

func NewListener(logger *zap.SugaredLogger, kafkaConsumer kf.Consumer, db service.DB) *Listener {
	return &Listener{
		logger:        logger,
		kafkaConsumer: kafkaConsumer,
		db:            db,
	}
}

func (e *Listener) Listen(ctx context.Context) {
	e.logger.Info("Listening for messages.")
	for {
		select {
		case <-ctx.Done():
			e.logger.Infow("Stopping listening messages")
			e.kafkaConsumer.Close()

			return
		default:
			e.handleMessage(ctx)
		}
	}
}

func (e *Listener) handleMessage(ctx context.Context) {
	commit, id, err := e.kafkaConsumer.ReadMessage(ctx)
	if err != nil {
		e.logger.Errorw("Unmarshal kafka message.", "err", err)
		return
	}

	status := "received"
	err = e.db.UpdateMessages(ctx, status, []int{id})
	if err != nil {
		e.logger.Infow("failed update messages", "err", err)
		return
	}
	fmt.Println(id, "-----------")

	err = commit(ctx)
	if err != nil {
		e.logger.Errorw("Commit kafka message.", "err", err)
	}
}
