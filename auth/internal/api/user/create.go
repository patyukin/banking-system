package user

import (
	"context"
	"github.com/patyukin/banking-system/auth/internal/converter"
	desc "github.com/patyukin/banking-system/auth/pkg/user_v1"
	"github.com/pkg/errors"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, errors.New("passwords do not match")
	}

	user := converter.ToUserFromDesc(req)
	id, err := i.userService.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}
