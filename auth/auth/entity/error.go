package entity

import (
	"github.com/regismelgaco/go-sdks/erring"
)

var (
	ErrExpiredToken         = erring.New("expired token", erring.WithType(erring.ErrUnauthorized))
	ErrInvalidToken         = erring.New("invalid token", erring.WithType(erring.ErrUnauthorized))
	ErrUsernameInUse        = erring.New("username in use", erring.WithType(erring.ErrBadRequest))
	ErrUserNotFound         = erring.New("user not found", erring.WithType(erring.ErrNotFound))
	ErrInvalidLoginInput    = erring.New("invalid login input", erring.WithType(erring.ErrBadRequest))
	ErrMissingAuthorization = erring.New("authorization header is missing", erring.WithType(erring.ErrBadRequest))
	ErrMissingClaimsCtx     = erring.New("authorization token claims are missing from context")
)
