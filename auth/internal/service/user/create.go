package user

import (
	"context"

	"github.com/patyukin/banking-system/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.UserInfo) (string, error) {
	var uuid string
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		uuid, errTx = s.userRepository.Create(ctx, info)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.userRepository.Get(ctx, uuid)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return uuid, nil
}
