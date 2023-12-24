package converter

import (
	"github.com/patyukin/banking-system/auth/internal/model"
	desc "github.com/patyukin/banking-system/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ToUserFromDesc(user *desc.CreateUserRequest) *model.User {
	return &model.User{
		Password:  user.Password,
		Info:      *ToUserInfoFromDesc(user.GetInfo()),
		CreatedAt: time.Now(),
	}
}

func ToUserFromService(user *model.User) *desc.User {
	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromService(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(time.Now()),
	}
}

func ToUserInfoFromService(info model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:  info.Name,
		Email: info.Email,
	}
}

func ToUserInfoFromDesc(info *desc.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
	}
}
