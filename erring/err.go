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
	list := []string{}

	if e.Name != "" {
		list = append(list, "name: "+e.Name)
	}
	if e.Description != "" {
		list = append(list, "description: "+e.Description)
	}
	if len(e.Payload) > 0 {
		list = append(list, fmt.Sprintf("payload: %v", e.Payload))
	}
	if e.TypeErr != nil {
		list = append(list, "type: "+e.TypeErr.Error())
	}
	if e.InternalErr != nil {
		list = append(list, "internal: "+e.InternalErr.Error())
	}

	return strings.Join(list, " | ")
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

func Wrap(err error) Err {
	e := Err{}.Wrap(err)

	return e
}

func (err Err) Err() error {
	return err
}

func (err Err) With(label string, v any) Err {
	if err.Payload == nil {
		err.Payload = make(map[string]any)
	}

	err.Payload[label] = v

	return err
}

func (b Err) Wrap(err error) Err {
	erringErr, ok := err.(Err)
	if !ok {
		return Err{
			InternalErr: err,
			Stack:       debug.Stack(),
		}
	}

	if isStackEnabled && len(b.Stack) == 0 {
		if len(erringErr.Stack) > 0 {
			b.Stack = erringErr.Stack
		} else {
			b.Stack = debug.Stack()
		}
	}

	if b.Name != "" {
		b.InternalErr = Err{Name: erringErr.Name}
	} else {
		b.Name = erringErr.Name
	}

	des := []string{}
	if b.Description != "" {
		des = append(des, b.Description)
	}
	if erringErr.Description != "" {
		des = append(des, erringErr.Description)
	}
	b.Description = strings.Join(des, " - ")

	if b.Payload == nil {
		b.Payload = map[string]any{}
	}
	if erringErr.Payload == nil {
		erringErr.Payload = map[string]any{}
	}
	for k, v := range erringErr.Payload {
		b.Payload[k] = v
	}

	if b.TypeErr == nil {
		b.TypeErr = erringErr.TypeErr
	}

	if b.InternalErr != nil {
		b.InternalErr = erringErr.InternalErr
	}

	return b
}

func (err Err) Describe(s string) Err {
	if err.Description != "" {
		err.Description += " - "
	}

	err.Description += s

	return err
}

func (err Err) ChangeType(t error) Err {
	err.TypeErr = t

	return err
}
