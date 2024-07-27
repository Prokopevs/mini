package main

import (
	"fmt"
	"os"
)

type envConfig struct {
	kafkaBroker  string
	kafkaTopic   string
	kafkaGroupID string
	pgConnString string
	httpAddr     string
}

func loadEnvConfig() (*envConfig, error) {
	const (
		provideEnvErrorMsg = `please provide "%s" environment variable`

		kafkaBroker     = "KAFKA_BROKER"
		kafkaTopic      = "KAFKA_TOPIC"
		kafkaGroupID    = "KAFKA_GROUP_ID"
		pgConnStringEnv = "PG_CONN"
		httpAddr        = "HTTP_ADDR"
	)

	var ok bool

	cfg := &envConfig{}

	cfg.kafkaBroker, ok = os.LookupEnv(kafkaBroker)
	if !ok {
		return nil, fmt.Errorf(provideEnvErrorMsg, kafkaBroker)
	}

	cfg.kafkaTopic, ok = os.LookupEnv(kafkaTopic)
	if !ok {
		return nil, fmt.Errorf(provideEnvErrorMsg, kafkaTopic)
	}

	cfg.kafkaGroupID, ok = os.LookupEnv(kafkaGroupID)
	if !ok {
		return nil, fmt.Errorf(provideEnvErrorMsg, kafkaGroupID)
	}

	cfg.pgConnString, ok = os.LookupEnv(pgConnStringEnv)
	if !ok {
		return nil, fmt.Errorf(provideEnvErrorMsg, pgConnStringEnv)
	}

	cfg.httpAddr, ok = os.LookupEnv(httpAddr)
	if !ok {
		return nil, fmt.Errorf(provideEnvErrorMsg, httpAddr)
	}

	return cfg, nil
}
