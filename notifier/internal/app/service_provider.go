package app

import (
	"context"
	"github.com/patyukin/banking-system/notifier/internal/api/notifier"
	"github.com/patyukin/banking-system/notifier/internal/queue/kafka"
	"log"

	"github.com/patyukin/banking-system/notifier/internal/config"
	"github.com/patyukin/banking-system/notifier/internal/config/env"
)

type serviceProvider struct {
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	notifierImpl *notifier.Implementation

	consumer kafka.KafkaConsumer
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) NotifierImpl(_ context.Context) *notifier.Implementation {
	if s.notifierImpl == nil {
		s.notifierImpl = notifier.NewImplementation()
	}

	return s.notifierImpl
}
