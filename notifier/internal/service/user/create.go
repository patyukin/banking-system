package user

import (
	"context"
	"github.com/patyukin/banking-system/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (s *serv) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64
	// Генерация хеша пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user.Password = string(hashedPassword)
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.Create(ctx, user)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
