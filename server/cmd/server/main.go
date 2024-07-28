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

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugaredLogger := logger.Sugar()

	dbConn, err := db.NewDatabase(context.Background(), sugaredLogger, cfg.pgConnString)
	if err != nil {
		return err
	}

	bufSize, err := strconv.Atoi(cfg.bufferSize)
	if err != nil {
		return err
	}

	kafkaProducer := kafka.NewProducer(cfg.kafkaTopic, []string{cfg.kafkaBroker}, cfg.kafkaGroupID)
	defer kafkaProducer.Close()
	kafkaConsumer := kafka.NewConsumer(cfg.kafkaTopic, []string{cfg.kafkaBroker}, cfg.kafkaGroupID)
	defer kafkaConsumer.Close()

	messageService := service.NewMessageService(dbConn, kafkaProducer, kafkaConsumer, sugaredLogger, bufSize)

	httpServer := handler.NewHTTP(cfg.httpAddr, sugaredLogger, messageService)

	sender := worker.NewSender(10*time.Second, messageService, sugaredLogger.With("pkg", "sender"))
	reader := worker.NewListener(sugaredLogger.With("pkg", "reader"), messageService)

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
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan
	cancel()

	wg.Wait()

	return nil
}

//	@title mini server
//	@version 1.0
//	@description This is mini server
//	@host mini.eridani.site
// 	@schemes http https
//	@BasePath /api/v1
func main() {
	err := run()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(exitCodeInitError)
	}
}
