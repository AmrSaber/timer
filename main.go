package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gosuri/uilive"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "missing duration")
		os.Exit(1)
	}

	var waitDuration time.Duration
	parsingErrs := make([]error, 0)
	for _, amount := range os.Args[1:] {
		duration, err := time.ParseDuration(amount)
		if err != nil {
			parsingErrs = append(parsingErrs, fmt.Errorf("invalid duration %q", amount))
			continue
		}

		waitDuration += duration
	}

	if len(parsingErrs) > 0 {
		parsingErr := errors.Join(parsingErrs...)
		fmt.Fprintf(os.Stderr, "error parsing duration(s):\n%s\n", parsingErr)
		os.Exit(1)
	}

	waitDuration = waitDuration.Round(time.Second)

	uiWriter := uilive.New()
	uiWriter.Start()

	now := time.Now()
	target := now.Add(waitDuration)
	for now.Before(target) {
		remaining := target.Sub(now)

		fmt.Fprintf(
			uiWriter,
			"%02d:%02d:%02d\n",
			int(remaining.Hours()),
			int(remaining.Minutes())%60,
			(int(remaining.Seconds())-int(remaining.Minutes())*60)%60,
		)

		time.Sleep(253 * time.Millisecond)
		now = time.Now()
	}

	fmt.Fprintln(uiWriter, "")
	uiWriter.Stop()

	// Move cursor 1 line up
	fmt.Print("\033[1A")
}
