package user

import (
	"context"
	"github.com/patyukin/banking-system/auth/internal/converter"
	desc "github.com/patyukin/banking-system/auth/pkg/user_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userObj, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: converter.ToUserFromService(userObj),
	}, nil
}
