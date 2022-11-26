package erring

import (
	"bytes"
	"encoding/base64"
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
	Description string
	Payload     map[string]any
	InternalErr error
}

type Err struct {
	//TODO Error code
	ErrBuilder
	Stack   []byte
	TypeErr error
}

func (e Err) Error() string {
	stack := []byte("")
	_, _ = base64.NewDecoder(base64.RawStdEncoding, bytes.NewReader(e.Stack)).Read(stack)

	return fmt.Sprintf(
		"name: %s\ndescription: %v\npayload: %v\ninternalErr: %v\ntypeErr: %v\n%s",
		e.Name,
		e.Description,
		e.Payload,
		e.InternalErr,
		e.TypeErr,
		stack,
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

func (b ErrBuilder) Wrap(err error) ErrBuilder {
	if e, ok := err.(Err); ok && e.Name != "" {
		b.Name = e.Name

		for k, v := range e.Payload {
			b.Payload[k] = v
		}

		return b
	}

	b.InternalErr = err

	return b
}

func (b ErrBuilder) Describe(s string) ErrBuilder {
	b.Description = s

	return b
}
