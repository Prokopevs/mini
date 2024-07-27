package service

import (
	"context"

	"github.com/Prokopevs/mini/server/internal/model"
)

func (s *MessageService) CreateMessage(ctx context.Context, message *model.MessageCreate) error {
	_, err := s.db.CreateMessage(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) GetMessages(ctx context.Context) ([]*model.Message, error) {
	m, err := s.db.GetMessages(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *MessageService) UpdateMessages(ctx context.Context, status string, ids []int) error {
	err := s.db.UpdateMessages(ctx, status, ids)
	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) GetMessagesEvent(ctx context.Context, limit int) ([]int, error) {
	m, err := s.db.GetMessagesEvent(ctx, limit)
	if err != nil {
		return nil, err
	}

	return m, nil
}