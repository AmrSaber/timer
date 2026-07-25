// Package common for common utils and helpers
package common

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

func Assert(cond bool, err any) {
	if !cond {
		panic(err)
	}
}

func FailOn(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleSigterm(ctx context.Context, fn func(os.Signal)) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

		select {
		case signal := <-signals:
			fn(signal)
		case <-ctx.Done():
		}
	}()
}

func ParseDuration() (time.Duration, error) {
	numberRegex := regexp.MustCompile(`^-?\d+(?:\.\d+)?$`)

	var waitDuration time.Duration
	parsingErrs := make([]error, 0)
	for _, amount := range os.Args[1:] {
		if numberRegex.MatchString(amount) {
			amount = amount + "s"
		}

		duration, err := time.ParseDuration(amount)
		if err != nil {
			parsingErrs = append(parsingErrs, fmt.Errorf("invalid duration %q", amount))
			continue
		}

		waitDuration += duration
	}

	if len(parsingErrs) > 0 {
		parsingErr := errors.Join(parsingErrs...)
		return 0, parsingErr
	}

	waitDuration = waitDuration.Round(time.Second)

	return waitDuration, nil
}
