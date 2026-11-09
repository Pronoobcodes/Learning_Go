package store

import (
	"context"

	"notes-api/internal/store/dbgen"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (dbgen.User, error)
	GetUserByEmail(ctx context.Context, email string) (dbgen.User, error)
	GetUserByID(ctx context.Context, id int64) (dbgen.User, error)
}

type Users struct {
	queries *dbgen.Queries
}

func NewUsers(queries *dbgen.Queries) *Users {
	return &Users{queries: queries}
}

func (u *Users) CreateUser(
	ctx context.Context,
	email, passwordHash string,
) (dbgen.User, error) {
	return u.queries.CreateUser(ctx, dbgen.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
}

func (u *Users) GetUserByEmail(
	ctx context.Context,
	email string,
) (dbgen.User, error) {
	return u.queries.GetUserByEmail(ctx, email)
}

func (u *Users) GetUserByID(
	ctx context.Context,
	id int64,
) (dbgen.User, error) {
	return u.queries.GetUserByID(ctx, id)
}
