package repository_test

import (
	"context"
	"testing"

	"github.com/regismelgaco/go-sdks/auth/auth/entity"
	"github.com/regismelgaco/go-sdks/auth/auth/gateway/postgres/repository"
	"github.com/regismelgaco/go-sdks/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_Repository_Insert(t *testing.T) {
	t.Parallel()

	type args struct {
		user entity.User
	}

	testCases := []struct {
		name    string
		payload []entity.User
		args
		want    entity.User
		wantErr error
	}{
		{
			name: "when insert new user expect user with updated id and external id",
			args: args{
				user: entity.User{
					UserName: "el_chavo",
					Secret:   []byte("8"),
				},
			},
			want: entity.User{
				ID:       1,
				UserName: "el_chavo",
				Secret:   []byte("8"),
			},
			wantErr: nil,
		},
		{
			name: "when insert user with used username expect used username error",
			payload: []entity.User{
				{
					UserName: "quico",
					Secret:   []byte("squared ball"),
				},
			},
			args: args{
				user: entity.User{
					UserName: "quico",
					Secret:   []byte("squared ball"),
				},
			},
			want: entity.User{
				UserName: "quico",
				Secret:   []byte("squared ball"),
			},
			wantErr: entity.ErrUsernameInUse,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pool := postgres.GetDB(t)

			SeedUsers(t, pool, tc.payload)

			r := repository.NewUserRepository(pool)

			err := r.Insert(context.Background(), &tc.args.user)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, tc.args.user)
		})
	}
}

func Test_Repository_SelectByUserName(t *testing.T) {
	t.Parallel()

	type args struct {
		username string
	}

	testCases := []struct {
		name    string
		payload []entity.User
		args
		want    entity.User
		wantErr error
	}{
		{
			name: "when username exists expect user",
			args: args{
				username: "el_chavo",
			},
			payload: []entity.User{
				{
					UserName: "el_chavo",
					Secret:   []byte("8"),
				},
			},
			want: entity.User{
				ID:       1,
				UserName: "el_chavo",
				Secret:   []byte("8"),
			},
			wantErr: nil,
		},
		{
			name: "when username does not exists expect error user not found",
			args: args{
				username: "quico",
			},
			want:    entity.User{},
			wantErr: entity.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pool := postgres.GetDB(t)

			SeedUsers(t, pool, tc.payload)

			r := repository.NewUserRepository(pool)

			got, err := r.SelectByUserName(context.Background(), tc.args.username)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}
