package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/patyukin/banking-system/auth/internal/model"
	desc "github.com/patyukin/banking-system/auth/pkg/user_v1"
)

func ToNoteFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Uuid:      user.UUID,
		Info:      ToNoteInfoFromService(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToNoteInfoFromService(info model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:  info.Name,
		Email: info.Email,
	}
}

func ToNoteInfoFromDesc(info *desc.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
	}
}
