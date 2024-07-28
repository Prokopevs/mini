package handler

import (
	"context"
	"errors"

	"net/http"

	"github.com/Prokopevs/mini/server/internal/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Service interface {
	CreateMessage(ctx context.Context, mess *model.MessageCreate) error
	GetMessages(ctx context.Context) ([]*model.Message, error)
	UpdateMessages(ctx context.Context, status string, ids []int) error
	GetMessagesEvent(ctx context.Context, limit int) ([]int, error)
}

type HTTP struct {
	innerServer *http.Server

	log     *zap.SugaredLogger
	service Service
}

func (h *HTTP) Run(ctx context.Context) {
	h.log.Infow("HTTP server starting.", "addr", h.innerServer.Addr)

	go func() {
		err := h.innerServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			h.log.Errorw("Listen and serve HTTP", "addr", h.innerServer.Addr, "err", err)
		}
	}()

	<-ctx.Done()

	h.log.Info("Graceful server shutdown.")
	h.innerServer.Shutdown(context.Background())
}

func NewHTTP(addr string, logger *zap.SugaredLogger, messageSvc Service) *HTTP {
	h := &HTTP{
		log:     logger,
		service: messageSvc,
	}

	r := gin.Default()
	h.setRoutes(r)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	h.innerServer = srv

	return h
}
