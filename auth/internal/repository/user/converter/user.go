package converter

import (
	"github.com/patyukin/banking-system/auth/internal/model"
	modelRepo "github.com/patyukin/banking-system/auth/internal/repository/user/model"
)

func ToNoteFromRepo(note *modelRepo.User) *model.User {
	return &model.User{
		UUID:      note.UUID,
		Info:      ToNoteInfoFromRepo(note.Info),
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func ToNoteInfoFromRepo(info modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
	}
}
