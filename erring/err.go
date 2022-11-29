package erring

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
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
	//TODO Error code
	name        string
	description string
	payload     map[string]any
	internalErr error
	typeErr     error
	stack       []byte
}

type Err struct {
	//TODO Error code
	Name        string
	Description string
	Payload     map[string]any
	InternalErr error
	TypeErr     error
	Stack       []byte
}

func (e Err) Error() string {
	return fmt.Sprintf(
		"name: %s\ndescription: %v\npayload: %v\ninternalErr: %v\ntypeErr: %v\n\n\n%s",
		e.Name,
		e.Description,
		e.Payload,
		e.InternalErr,
		e.TypeErr,
		e.Stack,
	)
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
	return ErrBuilder{}.Wrap(err)
}

func (b ErrBuilder) Build() error {
	if isStackEnabled && len(b.stack) > 0 {
		b.stack = debug.Stack()
	}

	return Err{
		Name:        b.name,
		Description: b.description,
		Payload:     b.payload,
		InternalErr: b.internalErr,
		TypeErr:     b.typeErr,
		Stack:       b.stack,
	}
}

func (b ErrBuilder) With(label string, v any) ErrBuilder {
	if b.payload == nil {
		b.payload = make(map[string]any)
	}

	b.payload[label] = v

	return b
}

func (b ErrBuilder) Wrap(err error) ErrBuilder {
	e, ok := err.(Err)
	if !ok {
		return ErrBuilder{internalErr: err}
	}

	if b.name != "" {
		b.internalErr = Err{Name: e.Name}
	} else {
		b.name = e.Name
	}

	e.Description = strings.Join([]string{b.description, e.Description}, " - ")

	if e.Payload == nil {
		e.Payload = map[string]any{}
	}
	if b.payload == nil {
		b.payload = map[string]any{}
	}
	for k, v := range e.Payload {
		b.payload[k] = v
	}

	if b.typeErr == nil {
		b.typeErr = e.TypeErr
	}

	if e.InternalErr != nil {
		b.internalErr = e.InternalErr
	}

	if e.Stack != nil {
		b.stack = e.Stack
	}

	return b
}

func (b ErrBuilder) Describe(s string) ErrBuilder {
	b.description = s

	return b
}

func (b ErrBuilder) ChangeType(err error) ErrBuilder {
	b.typeErr = err

	return b
}

func (b ErrBuilder) Err() Err {
	return Err{
		Name:        b.name,
		Description: b.description,
		Payload:     b.payload,
		InternalErr: b.internalErr,
		TypeErr:     b.typeErr,
		Stack:       b.stack,
	}
}
