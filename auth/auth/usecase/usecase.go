package usecase

import (
	"context"
	"time"

	"github.com/regismelgaco/go-sdks/auth/auth/entity"
	"github.com/regismelgaco/go-sdks/erring"
)

type Usecase interface {
	CreateUser(context.Context, CreateUserInput) (entity.User, error)
	Login(context.Context, LoginInput) (entity.Token, error)
	IsAuthorized(ctx context.Context, token entity.Token) (entity.TokenClaims, error)
}

type CreateUserInput struct {
	UserName string
	Pass     string
}

type LoginInput struct {
	UserName string
	Pass     string
}

type usecase struct {
	Encryptor
	Repository
}

func NewUsecase(encryptor Encryptor, repository Repository) Usecase {
	return usecase{Encryptor: encryptor, Repository: repository}
}

func (u usecase) CreateUser(ctx context.Context, input CreateUserInput) (entity.User, error) {
	secret, err := u.GenerateSecret(input.Pass)
	if err != nil {
		return entity.User{}, err
	}

	user := entity.User{
		UserName: input.UserName,
		Secret:   secret,
	}

	if err := u.Insert(ctx, &user); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u usecase) Login(ctx context.Context, input LoginInput) (entity.Token, error) {
	user, err := u.SelectByUserName(ctx, input.UserName)
	if err != nil {
		return "", erring.Wrap(entity.ErrInvalidLoginInput).Wrap(err)
	}

	if err := u.CheckPassword(user.Secret, input.Pass); err != nil {
		return "", erring.Wrap(entity.ErrInvalidLoginInput).Wrap(err)
	}

	const day = 24 * time.Hour

	claims := entity.TokenClaims{
		UserName:   input.UserName,
		Expiration: time.Now().Add(day),
	}

	token, err := u.GenerateSessionToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u usecase) IsAuthorized(ctx context.Context, token entity.Token) (entity.TokenClaims, error) {
	claims, err := u.GetClaimsFromToken(token)
	if err != nil {
		return entity.TokenClaims{}, erring.Wrap(err).Wrap(entity.ErrInvalidToken).Describe("failed to get claims from token")
	}

	if time.Now().After(claims.Expiration) {
		return entity.TokenClaims{}, erring.Wrap(err).Wrap(entity.ErrExpiredToken).Describe("expired session")
	}

	return claims, nil
}
