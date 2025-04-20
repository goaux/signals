package signals_test

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/goaux/signals"
)

func Example_wait() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Waiting for SIGINT or SIGTERM...")
	sig := signals.Wait(ctx, syscall.SIGINT, syscall.SIGTERM)

	if sig != nil {
		fmt.Printf("Received signal: %v\n", sig)
	} else {
		fmt.Println("Context canceled")
	}
	// Output:
	// Waiting for SIGINT or SIGTERM...
	// Context canceled
}

func TestWait(t *testing.T) {
	t.Run("Signal received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		go func() {
			time.Sleep(100 * time.Millisecond)
			syscall.Kill(os.Getpid(), syscall.SIGINT)
		}()

		sig := signals.Wait(ctx, syscall.SIGINT)
		if sig != syscall.SIGINT {
			t.Errorf("Expected SIGINT, got %v", sig)
		}
	})

	t.Run("Context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		sig := signals.Wait(ctx, syscall.SIGINT)
		if sig != nil {
			t.Errorf("Expected nil, got %v", sig)
		}
	})

	t.Run("Multiple signals", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		go func() {
			time.Sleep(100 * time.Millisecond)
			syscall.Kill(os.Getpid(), syscall.SIGTERM)
		}()

		sig := signals.Wait(ctx, syscall.SIGINT, syscall.SIGTERM)
		if sig != syscall.SIGTERM {
			t.Errorf("Expected SIGTERM, got %v", sig)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		sig := signals.Wait(ctx, syscall.SIGINT)
		if sig != nil {
			t.Errorf("Expected nil (timeout), got %v", sig)
		}
	})
}
