package repository

import (
	"context"

	"github.com/patyukin/banking-system/auth/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo) (string, error)
	Get(ctx context.Context, id string) (*model.User, error)
}
