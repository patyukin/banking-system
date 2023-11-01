package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/patyukin/banking-system/auth/internal/api/auth"
	"github.com/patyukin/banking-system/auth/internal/model"
	"github.com/patyukin/banking-system/auth/internal/service"
	serviceMocks "github.com/patyukin/banking-system/auth/internal/service/mocks"
	desc "github.com/patyukin/banking-system/auth/pkg/auth_v1"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func Test_SignIn(t *testing.T) {
	t.Parallel()
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.AuthRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		email        = gofakeit.Email()
		password     = gofakeit.Password(true, true, true, true, false, 10)
		accessToken  = gofakeit.UUID()
		refreshToken = gofakeit.UUID()

		hashPassword, _ = bcrypt.GenerateFromPassword([]byte(password), 2)

		serviceErr = fmt.Errorf("service error")

		req = &desc.AuthRequest{
			Email:    email,
			Password: password,
		}

		user = &model.User{
			Info:            model.UserInfo{},
			Roles:           nil,
			Password:        string(hashPassword),
			ConfirmPassword: string(hashPassword),
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		res = &desc.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.AuthResponse
		err             error
		authServiceMock authServiceMockFunc
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.SignInMock.Expect(ctx, user).Return(res, nil)

				return mock
			},
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetByEmailMock.Expect(ctx, email).Return(user, nil)

				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.SignInMock.Expect(ctx, user).Return(res, nil)

				return mock
			},
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetByEmailMock.Expect(ctx, email).Return(nil, serviceErr)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.authServiceMock(mc)
			userServiceMock := tt.userServiceMock(mc)
			api := auth.NewImplementation(authServiceMock, userServiceMock)

			response, err := api.SignIn(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}
