package signals

import (
	"context"
	"os"
	"os/signal"
)

// NotifyContext creates a context that is canceled when a signal is received,
// and returns the result of calling run with that context.
// The result of run is returned as is, regardless of whether it is an error.
func NotifyContext(
	ctx context.Context,
	run func(context.Context) error,
	signals ...os.Signal,
) error {
	ctx, cancel := signal.NotifyContext(ctx, signals...)
	defer cancel()
	return run(ctx)
}
