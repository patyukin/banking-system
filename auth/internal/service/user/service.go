package user

import (
	"github.com/patyukin/banking-system/auth/internal/client/db"
	"github.com/patyukin/banking-system/auth/internal/repository"
	"github.com/patyukin/banking-system/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
