package auth

import (
	"github.com/patyukin/banking-system/auth/internal/client/db"
	"github.com/patyukin/banking-system/auth/internal/repository"
	"github.com/patyukin/banking-system/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
	authRepository repository.AuthRepository
}

func NewService(authRepository repository.AuthRepository, userRepository repository.UserRepository, txManager db.TxManager) service.AuthService {
	return &serv{
		authRepository: authRepository,
		userRepository: userRepository,
		txManager:      txManager,
	}
}
