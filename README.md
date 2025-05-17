# signals
signals is a Go module that provides utilities for handling OS signals.

[![Go Reference](https://pkg.go.dev/badge/github.com/goaux/signals.svg)](https://pkg.go.dev/github.com/goaux/signals)
[![Go Report Card](https://goreportcard.com/badge/github.com/goaux/signals)](https://goreportcard.com/report/github.com/goaux/signals)

This package offers functionality to create contexts that can be canceled when specific
OS signals are received, making it easier to implement graceful shutdown or
interruption handling in applications.

The key feature is `NotifyContext`, which creates a context that is canceled when
one of the specified signals is received, and executes a function with that context.
The signal that caused cancellation can be retrieved using `FromContext`.

## Installation

To install the signals module, use the following command:

    go get github.com/goaux/signals

## Usage

### func NotifyContext and FromContext

`NotifyContext` creates a context that is canceled when a signal is received,
and runs a provided function with that context. The result of the function is
returned as is, regardless of whether it is an error.

`FromContext` can be used with the context to retrieve the signal that caused cancellation.
`context.Cause` can also be used to retrieve the signal within `Canceled`.

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/goaux/signals"
	"github.com/goaux/stacktrace/v2"
	"github.com/goaux/timer"
)

func main() {
	ctx := context.Background()
	err := signals.NotifyContext(ctx, run, os.Interrupt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", stacktrace.Format(err))
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	if err := timer.SleepCause(ctx, 7*time.Second); err != nil {
		if sig, ok := signals.FromContext(ctx); ok {
			fmt.Printf("Context was canceled by signal: %v\n", sig)
			return nil
		}
		return err
	}
	fmt.Println("hello")
	return nil
}
```

### func Wait

Here's a basic example of how to use the `Wait` function from the signals package:

```go
package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/goaux/signals"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	fmt.Println("Waiting for SIGINT or SIGTERM...")
	sig := signals.Wait(ctx, syscall.SIGINT, syscall.SIGTERM)

	if sig != nil {
		fmt.Printf("Received signal: %v\n", sig)
	} else {
		fmt.Println("Context canceled")
	}
}
```

This example waits for either SIGINT or SIGTERM for up to one minute. If a
signal is received, it prints the signal. If the context is canceled (in this
case, due to timeout), it prints "Context canceled".
