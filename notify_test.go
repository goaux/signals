package signals_test

import (
	"context"
	"fmt"
	"os"

	"github.com/goaux/signals"
)

func Example_notifyContext() {
	ctx := context.Background()
	err := signals.NotifyContext(ctx, run, os.Interrupt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	// Output:
	// hello
}

func run(ctx context.Context) error {
	fmt.Println("hello")
	return nil
}
