package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/AmrSaber/timer/src/common"
)

type TimerContextKey int

const signalKey TimerContextKey = iota

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "missing duration")
		os.Exit(1)
	}

	waitDuration, err := common.ParseDuration()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing duration(s):\n%s\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	common.HandleSigterm(ctx, func(signal os.Signal) {
		ctx = context.WithValue(ctx, signalKey, signal)
		cancel()
	})

	common.CountDown(ctx, waitDuration)

	if ctx.Err() != nil {
		if signal, ok := ctx.Value(signalKey).(syscall.Signal); ok {
			// If failed because of system signal: then exit with right code
			os.Exit(128 + int(signal))
		} else {
			// Otherwise: exit with generic code
			os.Exit(1)
		}
	}
}
