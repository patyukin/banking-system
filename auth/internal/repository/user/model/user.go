package model

import (
	"database/sql"
	"time"
)

type User struct {
	UUID      string
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserInfo struct {
	Email string
	Name  string
}
