package erring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	typeErr         = New("not found")
	internalErr     = errors.New("db conn lost")
	baseWithTypeErr = New("account not found", WithType(typeErr))
	baseErr         = New("account is blocked")
)

func Test_Is(t *testing.T) {
	t.Parallel()

	got := baseErr

	assert.ErrorIs(t, got, baseErr)
}

func Test_Type(t *testing.T) {
	t.Parallel()

	assert.ErrorIs(t, baseWithTypeErr, typeErr)
}

func Test_Wrap_Err(t *testing.T) {
	t.Parallel()

	DisableStack()

	got := Wrap(baseErr).With("number", 123).Build()

	assert.ErrorIs(t, got, baseErr)
	assert.JSONEq(t, `{
		"name": "account is blocked",
		"payload": {"number": 123}
	}`, got.Error())
}

func Test_Wrap_Internal(t *testing.T) {
	t.Parallel()

	DisableStack()

	got := Wrap(internalErr).With("number", 123).Build()

	assert.ErrorIs(t, got, internalErr)
	assert.JSONEq(t, `{
		"internal_error": "db conn lost",
		"payload": {"number": 123}
	}`, got.Error())
}
