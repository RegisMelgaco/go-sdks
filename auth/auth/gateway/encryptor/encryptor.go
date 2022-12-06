package encryptor

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/regismelgaco/go-sdks/auth/auth/entity"
	"github.com/regismelgaco/go-sdks/auth/auth/usecase"
	"github.com/regismelgaco/go-sdks/erring"
	"golang.org/x/crypto/bcrypt"
)

type encryptor struct {
	jwtSecret []byte
}

func NewEncryptor(jwtSecret string) usecase.Encryptor {
	if jwtSecret == "" {
		panic("jwtSecret is empty")
	}

	return encryptor{jwtSecret: []byte(jwtSecret)}
}

func (e encryptor) CheckPassword(secret []byte, pass string) error {
	err := bcrypt.CompareHashAndPassword(secret, []byte(pass))
	if err != nil {
		return erring.Wrap(err)
	}

	return nil
}

func (e encryptor) GenerateSecret(pass string) (secret []byte, err error) {
	s, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, erring.Wrap(err)
	}

	return s, nil
}

func (e encryptor) GenerateSessionToken(claims entity.TokenClaims) (entity.Token, error) {
	jwtClaims := jwt.MapClaims{
		"username":   claims.UserName,
		"expiration": claims.Expiration.Format(time.RFC3339),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	signedToken, err := t.SignedString(e.jwtSecret)
	if err != nil {
		return "", erring.Wrap(err)
	}

	return entity.Token(fmt.Sprintf("bearer %s", signedToken)), nil
}

func (e encryptor) GetClaimsFromToken(token entity.Token) (entity.TokenClaims, error) {
	s := strings.Split(string(token), " ")
	authenticationMethod, t := s[0], s[1]
	if strings.ToLower(authenticationMethod) != "bearer" {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken).
			With("got auth method", authenticationMethod)
	}

	parsed, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken).
				With("got sign method", t.Method.Alg())
		}

		return e.jwtSecret, nil
	})

	if err != nil {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken).Wrap(err)
	}

	if !parsed.Valid {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken).Describe("invalid token")
	}

	mapClaims := parsed.Claims.(jwt.MapClaims)

	v, ok := mapClaims["username"]
	if !ok {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken)
	}

	username, ok := v.(string)
	if err != nil {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken)
	}

	v, ok = mapClaims["expiration"]
	if !ok {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken)
	}

	expiration, err := time.Parse(time.RFC3339, v.(string))
	if err != nil {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken)
	}

	claims := entity.TokenClaims{
		UserName:   username,
		Expiration: expiration,
	}

	return claims, nil
}
