package user

import (
	"context"

	"github.com/patyukin/banking-system/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, uuid string) (*model.User, error) {
	user, err := s.userRepository.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return user, nil
}
