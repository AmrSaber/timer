package common

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/gosuri/uilive"
)

func CountDown(ctx context.Context, duration time.Duration) error {
	uiWriter := uilive.New()
	uiWriter.Start()

	defer func() {
		uiWriter.Stop()

		// Erase the trailing empty line so the timer vanishes cleanly
		fmt.Print("\033[1A")
	}()

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

		_, err := fmt.Fprintf(uiWriter, "%02d:%02d:%02d\n", hours, minutes, seconds)
		if err != nil {
			return err
		}

		// Odd interval to avoid syncing with second boundaries,
		// which would cause the display to hold the same value twice
		select {
		case <-time.After(253 * time.Millisecond):
		case <-ctx.Done():
			break loop
		}

		now = time.Now()
	}

	_, err := fmt.Fprintln(uiWriter, "")
	return err
}
