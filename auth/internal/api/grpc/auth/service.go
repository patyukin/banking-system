package auth

import (
	"github.com/patyukin/banking-system/auth/internal/service"
	desc "github.com/patyukin/banking-system/auth/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
	userService service.UserService
}

func NewImplementation(authService service.AuthService, userService service.UserService) *Implementation {
	return &Implementation{
		authService: authService,
		userService: userService,
	}
}
