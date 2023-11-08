package auth

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patyukin/banking-system/auth/internal/model"
	desc "github.com/patyukin/banking-system/auth/pkg/auth_v1"
)

const (
	TokenSignKey = "TOKEN_SIGN_KEY"
)

type Claims struct {
	Id int64 `json:"id"`
	jwt.RegisteredClaims
}

func (s *serv) SignIn(ctx context.Context, user *model.User) (*desc.AuthResponse, error) {
	// generate claims
	s.producer.SendMessage("ss", "sdsdsd")
	mySigningKey := []byte(os.Getenv(TokenSignKey))
	claims := Claims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	// generate access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}

	// generate refresh token
	refreshToken, errTx := s.authRepository.Create(ctx, user.ID)
	if errTx != nil {
		return nil, errTx
	}

	return &desc.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
