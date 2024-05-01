package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/gosuri/uilive"
)

var numberRegex = regexp.MustCompile(`^-?\d+$`)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "missing duration")
		os.Exit(1)
	}

	waitDuration, err := parseDuration()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing duration(s):\n%s\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handleSigterm(ctx, cancel)
	countDown(ctx, waitDuration)

	if ctx.Err() != nil {
		os.Exit(1)
	}
}

func parseDuration() (time.Duration, error) {
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

func handleSigterm(ctx context.Context, fn func()) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

		select {
		case <-signals:
			fn()
		case <-ctx.Done():
		}
	}()
}

func countDown(ctx context.Context, duration time.Duration) {
	uiWriter := uilive.New()
	uiWriter.Start()

	now := time.Now()
	target := now.Add(duration)

loop:
	for now.Before(target) {
		remaining := target.Sub(now)

		seconds := int(math.Ceil(remaining.Seconds()))
		minutes := seconds / 60
		hours := minutes / 60

		seconds %= 60
		minutes %= 60

		fmt.Fprintf(uiWriter, "%02d:%02d:%02d\n", hours, minutes, seconds)

		// Sleep or exit
		select {
		case <-time.After(253 * time.Millisecond):
		case <-ctx.Done():
			break loop
		}

		now = time.Now()
	}

	fmt.Fprintln(uiWriter, "")
	uiWriter.Stop()

	// Move cursor 1 line up
	fmt.Print("\033[1A")
}
