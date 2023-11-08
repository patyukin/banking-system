package app

import (
	"context"
	"log"

	"github.com/patyukin/banking-system/auth/internal/api/auth"
	"github.com/patyukin/banking-system/auth/internal/api/user"
	"github.com/patyukin/banking-system/auth/internal/queue/kafka"
	authRepository "github.com/patyukin/banking-system/auth/internal/repository/auth"

	"github.com/patyukin/banking-system/auth/internal/client/db"
	"github.com/patyukin/banking-system/auth/internal/client/db/pg"
	"github.com/patyukin/banking-system/auth/internal/client/db/transaction"
	"github.com/patyukin/banking-system/auth/internal/closer"
	"github.com/patyukin/banking-system/auth/internal/config"
	"github.com/patyukin/banking-system/auth/internal/config/env"
	"github.com/patyukin/banking-system/auth/internal/repository"
	userRepository "github.com/patyukin/banking-system/auth/internal/repository/user"
	"github.com/patyukin/banking-system/auth/internal/service"
	authService "github.com/patyukin/banking-system/auth/internal/service/auth"
	userService "github.com/patyukin/banking-system/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	dbClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository
	authRepository repository.AuthRepository

	userService service.UserService
	authService service.AuthService

	userImpl *user.Implementation
	authImpl *auth.Implementation
	producer *kafka.KafkaProducer
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
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

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.AuthRepository(ctx),
			s.UserRepository(ctx),
			s.TxManager(ctx),
			s.Producer(ctx),
		)
	}

	return s.authService
}

func (s *serviceProvider) Producer(_ context.Context) *kafka.KafkaProducer {
	return s.producer
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx), s.UserService(ctx))
	}

	return s.authImpl
}
