package errors

import (
	"fmt"
	"io"
)

func New(msg string) error {
	return &baseError{
		msg:   msg,
		stack: callers(),
	}
}

func Errorf(format string, args ...interface{}) error {
	return &baseError{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

type baseError struct {
	msg string
	stack
}

func (b *baseError) Error() string {
	return b.msg
}

func (b *baseError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, b.msg)
			b.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, b.msg)
	case 'q':
		fmt.Fprintf(s, "%q", b.msg)
	}
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}

type withStack struct {
	error
	stack
}

func (w *withStack) Cause() error {
	return w.error
}

func (w *withStack) Unwrap() error {
	return w.error
}

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

func WithMessage(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   msg,
	}
}

func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string {
	return w.msg + ": " + w.cause.Error()
}

func (w *withMessage) Cause() error {
	return w.cause
}

func (w *withMessage) Unwrap() error {
	return w.cause
}

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

func Cause(err error) error {
	for err != nil {
		cause, ok := err.(interface{ Cause() error })
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	err = &withMessage{
		cause: err,
		msg:   msg,
	}

	return &withStack{
		err,
		callers(),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}

	return &withStack{
		err,
		callers(),
	}
}
