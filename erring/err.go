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
	Err
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
	if isStackEnabled && len(b.Stack) > 0 {
		b.Stack = debug.Stack()
	}

	return b.Err
}

func (b ErrBuilder) With(label string, v any) ErrBuilder {
	if b.Payload == nil {
		b.Payload = make(map[string]any)
	}

	b.Payload[label] = v

	return b
}

func (b ErrBuilder) Wrap(err error) ErrBuilder {
	e, ok := err.(Err)
	if !ok {
		return ErrBuilder{Err{InternalErr: err}}
	}

	if b.Name != "" {
		b.InternalErr = Err{Name: e.Name}
	} else {
		b.Name = e.Name
	}

	e.Description = strings.Join([]string{b.Description, e.Description}, " - ")

	if e.Payload == nil {
		e.Payload = map[string]any{}
	}
	if b.Payload == nil {
		b.Payload = map[string]any{}
	}
	for k, v := range e.Payload {
		b.Payload[k] = v
	}

	if b.TypeErr == nil {
		b.TypeErr = e.TypeErr
	}

	if e.InternalErr != nil {
		b.InternalErr = e.InternalErr
	}

	if e.Stack != nil {
		b.Stack = e.Stack
	}

	return b
}

func (b ErrBuilder) Describe(s string) ErrBuilder {
	b.Description = s

	return b
}

func (b ErrBuilder) ChangeType(err error) ErrBuilder {
	b.TypeErr = err

	return b
}
