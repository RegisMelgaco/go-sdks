package erring

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
)

var (
	isStackEnabled = true
	stackToggleMux = sync.Mutex{}
)

func DisableStack() {
	stackToggleMux.Lock()
	isStackEnabled = false
	stackToggleMux.Unlock()
}

type ErrBuilder struct {
	Name        string
	Payload     map[string]any
	InternalErr error
}

type Err struct {
	ErrBuilder
	Stack   []byte
	TypeErr error
}

func (e Err) Error() string {
	type JSONErr struct {
		Name        string         `json:"name,omitempty"`
		Payload     map[string]any `json:"payload,omitempty"`
		InternalErr string         `json:"internal_error,omitempty"`
		Stack       []byte         `json:"stack,omitempty"`
		TypeErr     error          `json:"type,omitempty"`
	}

	var internalErr string
	if e.InternalErr != nil {
		internalErr = e.InternalErr.Error()
	}

	jErr := JSONErr{
		Name:        e.Name,
		Payload:     e.Payload,
		InternalErr: internalErr,
		Stack:       e.Stack,
		TypeErr:     e.TypeErr,
	}

	v, err := json.Marshal(jErr)
	if err != nil {
		return fmt.Sprintf(
			`{"name":%v,"payload":%v,"internalErr":%v,"typeErr":%v,"stack":%v}`,
			e.Name,
			e.Payload,
			e.InternalErr,
			e.TypeErr,
			e.Stack,
		)
	}

	return string(v)
}

func (e Err) Is(target error) bool {
	var t Err
	if errors.As(target, &t) {
		return e.Name == t.Name
	}

	return errors.Is(e.TypeErr, target) || errors.Is(e.InternalErr, target)
}

func (e Err) Unwrap() error {
	if e.TypeErr != nil {
		return e.TypeErr
	}

	return e.InternalErr
}

func New(name string, opts ...ErrOpt) error {
	var e Err
	e.Name = name

	for _, o := range opts {
		o(&e)
	}

	return e
}

type ErrOpt func(*Err)

func WithType(t error) func(*Err) {
	return func(e *Err) {
		e.TypeErr = t
	}
}

func Wrap(err error) ErrBuilder {
	b := ErrBuilder{}
	if e, ok := err.(Err); ok {
		b.Name = e.Name
		b.Payload = e.Payload

		return b
	}

	b.InternalErr = err

	return b
}

func (b ErrBuilder) Build() error {
	err := Err{}
	err.ErrBuilder = b

	if isStackEnabled {
		err.Stack = debug.Stack()
	}

	return err
}

func (b ErrBuilder) With(label string, v any) ErrBuilder {
	if b.Payload == nil {
		b.Payload = make(map[string]any)
	}

	b.Payload[label] = v

	return b
}
