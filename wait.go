// Package signals provides utilities for handling OS signals.
package signals

import (
	"context"
	"os"
	"os/signal"
)

// Wait waits for the specified OS signals or context cancellation.
// It returns the received signal or nil if the context is canceled.
//
// If no signals are provided, all incoming signals will be relayed.
// Otherwise, only the provided signals will be monitored.
//
// Multiple calls to Wait with the same signals are allowed and will work correctly:
// each call will receive copies of incoming signals independently.
func Wait(ctx context.Context, signals ...os.Signal) os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)
	select {
	case sig := <-ch:
		return sig
	case <-ctx.Done():
		return nil
	}
}
