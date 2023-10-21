package user

import "github.com/patyukin/banking-system/auth/internal/service"

type Handler struct {
	userService service.UserService
}

func New(userService service.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}
