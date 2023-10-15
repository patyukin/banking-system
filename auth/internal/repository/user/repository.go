package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/patyukin/banking-system/auth/internal/client/db"
	"github.com/patyukin/banking-system/auth/internal/model"
	"github.com/patyukin/banking-system/auth/internal/repository"
	"github.com/patyukin/banking-system/auth/internal/repository/user/converter"
	modelRepo "github.com/patyukin/banking-system/auth/internal/repository/user/model"
)

const (
	tableName = "users"

	uuidColumn      = "uuid"
	nameColumn      = "name"
	emailColumn     = "email"
	passwordColumn  = "password"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.UserInfo) (string, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn).
		Values(info.Name, info.Email).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	q := db.Query{
		Name:     "note_repository.Create",
		QueryRaw: query,
	}

	var uuid string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&uuid)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (r *repo) Get(ctx context.Context, uuid string) (*model.User, error) {
	builder := sq.Select(uuidColumn, nameColumn, emailColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{uuidColumn: uuid}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "note_repository.Get",
		QueryRaw: query,
	}

	var note modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).
		Scan(&note.UUID, &note.Info.Name, &note.Info.Email, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}
