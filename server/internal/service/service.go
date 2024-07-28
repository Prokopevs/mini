package service

import (
	"context"
	"github.com/Prokopevs/mini/server/internal/model"
)

type DB interface {
	CreateMessage(ctx context.Context, mess *model.MessageCreate) (int, error)
	GetMessages(ctx context.Context) ([]*model.Message, error)
	UpdateMessages(ctx context.Context, status string, ids []int) (error)
	GetMessagesEvent(ctx context.Context, limit int) ([]int, error)
	RunTx(ctx context.Context, f func(context.Context) error) error
}

type MessageService struct {
	db DB
}

func NewMessageService(db DB) *MessageService {
	return &MessageService{
		db: db,
	}
}
