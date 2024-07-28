package db

import (
	"context"
	"github.com/Prokopevs/mini/server/internal/model"
	"github.com/jmoiron/sqlx"
)

func (r *database) CreateMessage(ctx context.Context, mess *model.MessageCreate) (int, error) {
	const query = `INSERT INTO messages(message) VALUES ($1) RETURNING id`

	var id int

	err := r.GetExtContext(ctx).QueryRowxContext(ctx, query, mess.Message).Scan(&id)
	if err != nil {
		return 0, err
	}
	r.logger.Infow("Successfully added to db message with id", "id", id)
	return id, nil
}

func (r *database) GetMessages(ctx context.Context) ([]*model.Message, error) {
	const q = "SELECT * FROM messages ORDER BY id ASC"

	m := []*model.Message{}
	
	err := sqlx.SelectContext(ctx, r.GetQueryerContext(ctx), &m, q)
	if err != nil {
		return nil, err
	}

	return m, err
}

func (r *database) UpdateMessages(ctx context.Context, status string, ids []int) error {
	query, args, err := sqlx.In(`UPDATE messages SET status = ? WHERE id IN (?)`, status, ids)
    if err != nil {
        return err
    }

	query = r.GetExtContext(ctx).Rebind(query)
    _, err = r.GetExtContext(ctx).ExecContext(ctx, query, args...)
    if err != nil {
        return err
    }

    return nil
}

func (r *database) GetMessagesEvent(ctx context.Context, limit int) ([]int, error) {
	const q = "SELECT id FROM messages WHERE status = 'idle' ORDER BY id ASC LIMIT $1"

	m := []int{}

	err := sqlx.SelectContext(ctx, r.GetQueryerContext(ctx), &m, q, limit)
	if err != nil {
		return nil, err
	}

	return m, err
}