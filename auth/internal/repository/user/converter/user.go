package converter

import (
	"database/sql"
	"github.com/patyukin/banking-system/auth/internal/model"
	modelRepo "github.com/patyukin/banking-system/auth/internal/repository/user/model"
	"time"
)

func ToUserFromRepo(user *modelRepo.User) *model.User {
	var nullTime sql.NullTime
	var t time.Time

	if nullTime.Valid {
		t = nullTime.Time
	} else {
		t = time.Time{}
	}

	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: t,
	}
}

func ToUserInfoFromRepo(info modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
	}
}
