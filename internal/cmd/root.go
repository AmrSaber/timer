// Package cmd contains all the commands used.
package cmd

import (
	"context"
	"os"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/AmrSaber/timer/internal/common"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "timer <duration>",
	Short: "Countdown timer",
	Long: `Countdown tool that waits for a specified duration and exits when the time is up.
Arguments can be in the form "1h30m10s", "2m", "30s", ...
Duration is rounded to the nearest second.`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse duration
		var waitDuration time.Duration
		{
			value := args[0]

			// If value is a number
			if _, err := strconv.Atoi(value); err == nil {
				value = value + "s"
			}

			var err error
			waitDuration, err = time.ParseDuration(value)
			common.Assert(err == nil, "Error parsing duration: %v", err)

			waitDuration = waitDuration.Round(time.Second)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var exitSignal atomic.Pointer[syscall.Signal]
		common.HandleSigterm(ctx, func(signal syscall.Signal) {
			exitSignal.Store(&signal)
			cancel()
		})

		err := common.CountDown(ctx, waitDuration)
		common.FailOn(err)

		if ctx.Err() != nil {
			if signal := exitSignal.Load(); signal != nil {
				// If failed because of system signal: then exit with right code
				os.Exit(128 + int(*signal))
			} else {
				// Otherwise: exit with generic code
				os.Exit(1)
			}
		}
	},
}

// Execute sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Set version after it's been potentially injected in main.go
	rootCmd.Version = getVersion()

	err := rootCmd.Execute()
	common.FailOn(err)
}

func getVersion() string {
	version := common.GetVersion()
	if version == "" {
		return "??"
	}

	return version
}
