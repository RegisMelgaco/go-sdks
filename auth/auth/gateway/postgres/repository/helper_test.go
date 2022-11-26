package repository_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/regismelgaco/go-sdks/auth/auth/entity"
	"github.com/stretchr/testify/require"
)

func SeedUsers(t *testing.T, p *pgxpool.Pool, users []entity.User) {
	if users == nil {
		return
	}

	for _, u := range users {
		const sql = `
			insert into auth_user (username, secret)
			values ($1, $2)
			returning (id);
		`
		err := p.QueryRow(context.Background(), sql, u.UserName, u.Secret).Scan(&u.ID)
		require.NoError(t, err)
	}
}
