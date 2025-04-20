# signals
signals is a Go module that provides utilities for handling OS signals.

[![Go Reference](https://pkg.go.dev/badge/github.com/goaux/signals.svg)](https://pkg.go.dev/github.com/goaux/signals)
[![Go Report Card](https://goreportcard.com/badge/github.com/goaux/signals)](https://goreportcard.com/report/github.com/goaux/signals)

## Installation

To install the signals module, use the following command:

    go get github.com/goaux/signals

## Usage

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

### func NotifyContext

`NotifyContext` creates a context that is canceled when a signal is received,
and runs a provided function with that context. The result of the function is
returned as is, regardless of whether it is an error.

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/goaux/signals"
	"github.com/goaux/stacktrace/v2"
)

func main() {
	ctx := context.Background()
	err := signals.NotifyContext(ctx, run, os.Interrupt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", stacktrace.Format(err))
		os.Exit(1)
	}
	// Output:
	// hello
}

func run(ctx context.Context) error {
	fmt.Println("hello")
	return nil
}
```

## API

### func Wait

```go
func Wait(ctx context.Context, signals ...os.Signal) os.Signal
```

`Wait` waits for the specified OS signals or context cancellation. It returns
the received signal or nil if the context is canceled.
If no signals are provided, all incoming signals will be relayed. Otherwise,
only the provided signals will be monitored.
Multiple calls to Wait with the same signals are allowed and will work
correctly: each call will receive copies of incoming signals independently.

### func NotifyContext

```go
func NotifyContext(
	ctx context.Context,
	run func(context.Context) error,
	signals ...os.Signal,
) error {
```

`NotifyContext` creates a context that is canceled when a signal is received,
and returns the result of calling `run` with that context. The result of `run`
is returned as is, regardless of whether it is an error.
