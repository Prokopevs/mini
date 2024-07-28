package worker

import (
	"context"
	"time"

	"github.com/Prokopevs/mini/server/internal/service"
	"go.uber.org/zap"
)

type Sender struct {
	timeout        time.Duration
	logger         *zap.SugaredLogger
	messageService *service.MessageService
}

func NewSender(timeout time.Duration, messageService *service.MessageService, logger *zap.SugaredLogger) *Sender {
	return &Sender{
		timeout:        timeout,
		logger:         logger,
		messageService: messageService,
	}
}

func (t *Sender) RunProcessingMessages(ctx context.Context) {
	ticker := time.NewTicker(t.timeout)

	t.logger.Info("Run process messages")

	for {
		t.messageService.ProcessNewMessages(ctx)

		select {
		case <-ctx.Done():
			t.logger.Info("Stopping process messages")
			ticker.Stop()
			return
		case <-ticker.C:
		}
	}
}
