package usecase

import "github.com/regismelgaco/go-sdks/auth/auth/entity"

type Encryptor interface {
	GenerateSecret(pass string) ([]byte, error)
	CheckPassword(secret []byte, password string) error
	GenerateSessionToken(entity.TokenClaims) (entity.Token, error)
	GetClaimsFromToken(entity.Token) (entity.TokenClaims, error)
}
