package usecase

import (
	"context"

	"github.com/regismelgaco/go-sdks/auth/auth/entity"
)

type Repository interface {
	Insert(context.Context, *entity.User) error
	SelectByUserName(ctx context.Context, username string) (entity.User, error)
}
