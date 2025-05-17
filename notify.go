// Package signals provides utilities for handling OS signals in Go applications.
//
// This package offers functionality to create contexts that can be canceled when specific
// OS signals are received, making it easier to implement graceful shutdown or
// interruption handling in applications.
//
// The key feature is [NotifyContext], which creates a context that is canceled when
// one of the specified signals is received, and executes a function with that context.
// The signal that caused cancellation can be retrieved using [FromContext].
package signals

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
)

// NotifyContext creates a context that is canceled when one of the specified signals is received,
// and returns the result of calling run with that context.
// The result of run is returned as is, regardless of whether it is an error.
//
// The function takes a parent context, a function to run with the created context, and
// a list of signals to watch for. When one of these signals is received, the context
// is canceled with a Canceled error that wraps the signal.
//
// [FromContext] can be used with the context to retrieve the signal that caused cancellation.
// [context.Cause] can also be used to retrieve the signal within [Canceled].
func NotifyContext(
	parent context.Context,
	run func(context.Context) error,
	signals ...os.Signal,
) error {
	ctx, cancel := context.WithCancelCause(parent)
	defer cancel(nil)
	c := &signalCtx{
		Context: ctx,
		signals: signals,
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)
	if ctx.Err() == nil {
		go func() {
			select {
			case sig := <-ch:
				cancel(Canceled{signal: sig})
			case <-ctx.Done():
			}
		}()
	}
	return run(c)
}

// FromContext retrieves the signal that caused a context to be canceled.
//
// If the context was canceled by [NotifyContext] due to receiving a signal,
// FromContext returns that signal and true. Otherwise, it returns nil and false.
//
// Example:
//
//	if sig, ok := signals.FromContext(ctx); ok {
//		fmt.Printf("Context was canceled by signal: %v\n", sig)
//	}
func FromContext(ctx context.Context) (signal os.Signal, ok bool) {
	var err Canceled
	if errors.As(context.Cause(ctx), &err) {
		return err.Signal(), true
	}
	return
}

// Canceled is an error that wraps context.Canceled and includes the OS signal
// that caused the cancellation.
//
// It implements the error interface and provides methods to access the underlying
// signal and unwrap to context.Canceled.
type Canceled struct {
	signal os.Signal
}

// Error returns the error message, which includes the context.Canceled error
// and the string representation of the signal.
func (s Canceled) Error() string {
	return fmt.Sprintf("%v (%s)", context.Canceled, s.signal.String())
}

// Unwrap returns the underlying context.Canceled error.
func (s Canceled) Unwrap() error {
	return context.Canceled
}

// Signal returns the OS signal that caused the cancellation.
func (s Canceled) Signal() os.Signal {
	return s.signal
}

type signalCtx struct {
	context.Context
	signals []os.Signal
}

type stringer interface {
	String() string
}

func (c *signalCtx) String() string {
	var buf []byte
	// We know that the type of c.Context is context.cancelCtx, and we know that the
	// String method of cancelCtx returns a string that ends with ".WithCancel".
	name := c.Context.(stringer).String()
	name = name[:len(name)-len(".WithCancel")]
	buf = append(buf, "signals.NotifyContext("+name...)
	if len(c.signals) != 0 {
		buf = append(buf, ", ["...)
		for i, s := range c.signals {
			buf = append(buf, s.String()...)
			if i != len(c.signals)-1 {
				buf = append(buf, ' ')
			}
		}
		buf = append(buf, ']')
	}
	buf = append(buf, ')')
	return string(buf)
}
