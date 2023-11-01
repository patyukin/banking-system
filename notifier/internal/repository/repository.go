package repository

import (
	"context"

	"github.com/patyukin/banking-system/auth/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type AuthRepository interface {
	Create(ctx context.Context, userId int64) (string, error)
}
