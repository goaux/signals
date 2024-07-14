# signals
signals is a Go module that provides utilities for handling OS signals.

## Installation

To install the signals module, use the following command:

    go get github.com/goaux/signals

## Usage

Here's a basic example of how to use the `Wait` function from the signals package:

```go
package main

import (
    "context"
    "fmt"
    "os"
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
