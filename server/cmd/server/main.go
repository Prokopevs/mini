package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Prokopevs/mini/server/db"
	"github.com/Prokopevs/mini/server/internal/kafka"
	"github.com/Prokopevs/mini/server/internal/service"
	"github.com/Prokopevs/mini/server/internal/transport/message/handler"
	"github.com/Prokopevs/mini/server/internal/worker"
	"go.uber.org/zap"
)

const (
	exitCodeInitError = 2
)

func run() error {
	cfg, err := loadEnvConfig()
	if err != nil {
		return err
	}

	dbConn, err := db.NewDatabase(context.Background(), cfg.pgConnString)
	if err != nil {
		return err
	}

	kafkaProducer := kafka.NewProducer(cfg.kafkaTopic, []string{cfg.kafkaBroker}, cfg.kafkaGroupID)
	kafkaConsumer := kafka.NewConsumer(cfg.kafkaTopic, []string{cfg.kafkaBroker}, cfg.kafkaGroupID)

	messageSvc := service.NewMessageService(dbConn)

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugaredLogger := logger.Sugar()

	httpServer := handler.NewHTTP(cfg.httpAddr, sugaredLogger, messageSvc)

	buffSize, err := strconv.Atoi(cfg.bufferSize)
	if err != nil {
		return err
	}
	sender := worker.NewSender(10*time.Second, dbConn, sugaredLogger, kafkaProducer, buffSize)
	reader := worker.NewListener(sugaredLogger, kafkaConsumer, dbConn)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(ctx context.Context) {
		httpServer.Run(ctx)
		wg.Done()
	}(ctx)

	wg.Add(1)
	go func(ctx context.Context) {
		sender.RunProcessingMessages(ctx)
		wg.Done()
	}(ctx)

	wg.Add(1)
	go func(ctx context.Context) {
		reader.Listen(ctx)
		wg.Done()
	}(ctx)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-termChan
	cancel()

	wg.Wait()

	return nil
}

//  @title mini server
//  @version 1.0
//	@description This is mini server
//	@host localhost:5555
//	@BasePath /api/v1
func main() {
	err := run()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(exitCodeInitError)
	}
}