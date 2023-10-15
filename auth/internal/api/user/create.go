package user

import (
	"context"
	"log"

	"github.com/patyukin/banking-system/auth/internal/converter"
	desc "github.com/patyukin/banking-system/auth/pkg/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	uuid, err := i.userService.Create(ctx, converter.ToNoteInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %s", uuid)

	return &desc.CreateUserResponse{
		Uuid: uuid,
	}, nil
}
