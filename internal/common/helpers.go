// Package common for common utils and helpers
package common

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Assert(cond bool, message string, args ...any) {
	if !cond {
		fmt.Fprintf(os.Stderr, message, args...)
		os.Exit(1)
	}
}

func FailOn(err error) {
	Assert(err == nil, "%v", err)
}

func HandleSigterm(ctx context.Context, fn func(syscall.Signal)) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

		select {
		case signal := <-signals:
			fn(signal.(syscall.Signal))
		case <-ctx.Done():
		}
	}()
}
