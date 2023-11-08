package auth

import (
	"context"

	desc "github.com/patyukin/banking-system/auth/pkg/auth_v1"
	"golang.org/x/crypto/bcrypt"
)

func (i *Implementation) SignIn(ctx context.Context, req *desc.AuthRequest) (*desc.AuthResponse, error) {
	user, err := i.userService.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetPassword()))
	if err != nil {
		return nil, err
	}

	return i.authService.SignIn(ctx, user)
}
