package erring

import "errors"

var (
	ErrApi = errors.New("api error")

	ErrUnauthorized = New("unauthorized", WithType(ErrApi))
)
