package auth

import (
	"github.com/patyukin/banking-system/auth/internal/client/db"
	"github.com/patyukin/banking-system/auth/internal/queue/kafka"
	"github.com/patyukin/banking-system/auth/internal/repository"
	"github.com/patyukin/banking-system/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
	authRepository repository.AuthRepository
	producer       *kafka.KafkaProducer
}

func NewService(authRepository repository.AuthRepository, userRepository repository.UserRepository, txManager db.TxManager, producer *kafka.KafkaProducer) service.AuthService {
	return &serv{
		authRepository: authRepository,
		userRepository: userRepository,
		txManager:      txManager,
		producer:       producer,
	}
}
