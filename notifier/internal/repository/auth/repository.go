package auth

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/patyukin/banking-system/auth/internal/client/db"
	"github.com/patyukin/banking-system/auth/internal/repository"
	"time"
)

const (
	tableName = "refresh_tokens"

	idColumn           = "id"
	refreshTokenColumn = "refresh_token"
	userIdColumn       = "user_id"
	expiresAtColumn    = "expired_at"
	createdAtColumn    = "created_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, userId int64) (string, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(userIdColumn, expiresAtColumn).
		Values(userId, time.Now().Add(24*7*time.Hour)).
		Suffix("RETURNING token")

	query, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var token string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}
