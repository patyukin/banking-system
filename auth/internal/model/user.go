package model

import (
	"time"
)

type User struct {
	ID              int64
	Info            UserInfo
	Roles           []Role
	Password        string
	ConfirmPassword string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UserInfo struct {
	Name  string
	Email string
}

type Role struct {
	ID   int32
	Name string
}
