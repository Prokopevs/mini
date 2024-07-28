package service

import (
	"context"

	"github.com/Prokopevs/mini/server/internal/model"
	"go.uber.org/zap"
)

type txRunner interface {
	RunTx(ctx context.Context, f func(context.Context) error) error
}

type messagesRepo interface {
	txRunner

	CreateMessage(ctx context.Context, mess *model.MessageCreate) (int, error)
	GetMessages(ctx context.Context) ([]*model.Message, error)
	UpdateMessages(ctx context.Context, status string, ids []int) error
	GetMessagesEvent(ctx context.Context, limit int) ([]int, error)
}

// messagesNotifier - notify about message events.
type messagesNotifier interface {
	NotifyMessagesCreated(context.Context, []int) error
}

// messagesFetcher - fetch latest message events.
type messagesFetcher interface {
	FetchMessage(context.Context) (func(context.Context) error, int, error)
}

type MessageService struct {
	db                messagesRepo
	limit             int
	messagingNotifier messagesNotifier
	messagingFetcher  messagesFetcher
	logger            *zap.SugaredLogger
}

func NewMessageService(messagesRepo messagesRepo, messagingNotifier messagesNotifier, messagingFetcher messagesFetcher,
	logger *zap.SugaredLogger, limit int) *MessageService {
	return &MessageService{
		db:                messagesRepo,
		limit:             limit,
		messagingNotifier: messagingNotifier,
		messagingFetcher:  messagingFetcher,
		logger:            logger,
	}
}
