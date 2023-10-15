package user

import (
	"context"
	"log"

	"github.com/patyukin/banking-system/auth/internal/converter"
	desc "github.com/patyukin/banking-system/auth/pkg/user_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userObj, err := i.userService.Get(ctx, req.GetUuid())
	if err != nil {
		return nil, err
	}

	log.Printf(
		"uuid: %s, name: %s, email: %s, created_at: %v, updated_at: %v\n",
		userObj.UUID,
		userObj.Info.Name,
		userObj.Info.Email,
		userObj.CreatedAt,
		userObj.UpdatedAt,
	)

	return &desc.GetResponse{
		User: converter.ToNoteFromService(userObj),
	}, nil
}
