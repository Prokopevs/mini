package worker

import (
	"context"

	"github.com/Prokopevs/mini/server/internal/service"
	"go.uber.org/zap"
)

type Listener struct {
	logger         *zap.SugaredLogger
	messageService *service.MessageService
}

func NewListener(logger *zap.SugaredLogger, messageService *service.MessageService) *Listener {
	return &Listener{
		logger:         logger,
		messageService: messageService,
	}
}

func (e *Listener) Listen(ctx context.Context) {
	e.logger.Infow("Listening for messages.")
	for {
		select {
		case <-ctx.Done():
			e.logger.Infow("Stopping listening messages")

			return
		default:
			e.messageService.HandleQueryMessage(ctx)
		}
	}
}
