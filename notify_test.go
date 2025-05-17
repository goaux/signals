package signals_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/goaux/signals"
)

func Example_notifyContext() {
	ctx := context.Background()
	err := signals.NotifyContext(ctx, runN, syscall.SIGINT, syscall.SIGTERM)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	// Output:
	// signals.NotifyContext(context.Background, [interrupt terminated])
	// context canceled (interrupt)
	// true
	// interrupt true
}

func runN(ctx context.Context) error {
	fmt.Println(ctx)
	syscall.Kill(os.Getpid(), syscall.SIGINT)
	<-ctx.Done()
	err := context.Cause(ctx)
	fmt.Println(err.Error())
	fmt.Println(errors.Is(err, context.Canceled))
	fmt.Println(signals.FromContext(ctx))
	return nil
}
