package encryptor

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/regismelgaco/go-sdks/auth/auth/entity"
	"github.com/regismelgaco/go-sdks/auth/auth/usecase"
	"github.com/regismelgaco/go-sdks/erring"
	"golang.org/x/crypto/bcrypt"
)

type encryptor struct {
	jwtSecret []byte
}

func NewEncryptor(jwtSecret string) usecase.Encryptor {
	return encryptor{jwtSecret: []byte(jwtSecret)}
}

func (e encryptor) CheckPassword(secret []byte, pass string) error {
	err := bcrypt.CompareHashAndPassword(secret, []byte(pass))
	if err != nil {
		return erring.Wrap(err).Build()
	}

	return nil
}

func (e encryptor) GenerateSecret(pass string) (secret []byte, err error) {
	s, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, erring.Wrap(err).Build()
	}

	return s, nil
}

func (e encryptor) GenerateSessionToken(claims entity.TokenClaims) (entity.Token, error) {
	jwtClaims := jwt.MapClaims{
		"username":   claims.UserName,
		"expiration": claims.Expiration,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	signedToken, err := t.SignedString(e.jwtSecret)
	if err != nil {
		return "", erring.Wrap(err).Build()
	}

	return entity.Token(fmt.Sprintf("bearer %s", signedToken)), nil
}

func (e encryptor) GetClaimsFromToken(token entity.Token) (entity.TokenClaims, error) {
	s := strings.Split(string(token), " ")
	authenticationMethod, t := s[0], s[1]
	if strings.ToLower(authenticationMethod) != "bearer" {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken).Build()
	}

	parsed, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken).Build()
		}
		return e.jwtSecret, nil
	})

	if err != nil || !parsed.Valid {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken).Build()
	}

	mapClaims := parsed.Claims.(jwt.MapClaims)

	v, ok := mapClaims["username"]
	if !ok {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken).Build()
	}

	username, ok := v.(string)
	if err != nil {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken).Build()
	}

	v, ok = mapClaims["expiration"]
	if !ok {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken).Build()
	}

	expiration, err := time.Parse(time.RFC3339, v.(string))
	if err != nil {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrInvalidToken).Build()
	}

	claims := entity.TokenClaims{
		UserName:   username,
		Expiration: expiration,
	}

	return claims, nil
}
