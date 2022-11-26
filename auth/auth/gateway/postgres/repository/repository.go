package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/regismelgaco/go-sdks/auth/auth/entity"
	"github.com/regismelgaco/go-sdks/auth/auth/usecase"
	"github.com/regismelgaco/go-sdks/erring"
)

type repo struct {
	p *pgxpool.Pool
}

func NewUserRepository(p *pgxpool.Pool) usecase.Repository {
	return repo{p}
}

func (r repo) Insert(ctx context.Context, user *entity.User) error {
	const query = `
		insert into auth_user (username, secret)
		values ($1, $2) 
		returning (id);
	`

	err := r.p.QueryRow(ctx, query, user.UserName, user.Secret).Scan(&user.ID)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr); pgErr.Code == pgerrcode.UniqueViolation {
			return erring.Wrap(entity.ErrUsernameInUse).Build()
		}

		return erring.Wrap(err).Build()
	}

	return nil
}

func (r repo) SelectByUserName(ctx context.Context, username string) (entity.User, error) {
	const query = `
		select id, username, secret
		from auth_user
		where username = $1
		limit 1;
	`

	var user entity.User
	err := r.p.QueryRow(ctx, query, username).Scan(&user.ID, &user.UserName, &user.Secret)
	if err != nil {
		errBuilder := erring.Wrap(err).With("username", username)
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errBuilder.Wrap(entity.ErrUserNotFound).Build()
		}

		return entity.User{}, erring.Wrap(err).Describe("failed to select user").Build()
	}

	return user, nil
}
