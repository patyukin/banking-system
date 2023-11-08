package service

import (
	"context"

	desc "github.com/patyukin/banking-system/auth/pkg/auth_v1"

	"github.com/patyukin/banking-system/auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, info *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type AuthService interface {
	SignIn(ctx context.Context, params *model.User) (*desc.AuthResponse, error)
}
