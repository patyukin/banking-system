package model

import "time"

type RefreshToken struct {
	ID        int64
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
}

type AuthUser struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
