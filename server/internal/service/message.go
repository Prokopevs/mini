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

func (s *MessageService) ProcessNewMessages(ctx context.Context) {
	err := s.db.RunTx(ctx, func(ctx context.Context) error {
		// get messages with status='idle'
		ids, err := s.db.GetMessagesEvent(ctx, s.limit)
		if err != nil {
			s.logger.Errorw("failed to get messages", "err", err)
			return err
		}

		if len(ids) == 0 {
			return nil
		}

		// send to kafka
		err = s.messagingNotifier.NotifyMessagesCreated(ctx, ids)
		if err != nil {
			s.logger.Errorw("write to kafka", "err", err)
			return err
		}
		s.logger.Infow("Successfully sended messages to kafka", "ids", ids)

		// mark messages as sended
		status := "send"
		err = s.db.UpdateMessages(ctx, status, ids)
		if err != nil {
			s.logger.Errorw("failed update messages", "err", err)
			return err
		}
		s.logger.Infow("Changed statuses to `send` in db for", "ids", ids)

		return nil
	})
	if err != nil {
		s.logger.Errorw("tx processMessages error", "err", err)
	}
}

func (s *MessageService) HandleQueryMessage(ctx context.Context) {
	ack, id, err := s.messagingFetcher.FetchMessage(ctx)
	if err != nil {
		s.logger.Errorw("Unmarshal kafka message.", "err", err)
		return
	}

	status := "received"
	err = s.db.UpdateMessages(ctx, status, []int{id})
	if err != nil {
		s.logger.Errorw("failed update messages", "err", err)
		return
	}

	err = ack(ctx)
	if err != nil {
		s.logger.Errorw("Commit kafka message.", "err", err)
	}
	s.logger.Infow("Successfully readed from kafka", "id", id)
	s.logger.Infow("Changed statuses to `received` in db for", "id", id)
}
