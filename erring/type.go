package erring

var (
	ErrUnauthorized = New("unauthorized")
	ErrBadRequest   = New("bad request")
	ErrNotFound     = New("not found")
)
