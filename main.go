package main

import (
	"github.com/AmrSaber/timer/internal/cmd"
	"github.com/AmrSaber/timer/internal/common"
)

var version string

func main() {
	// Set version number if it's loaded from build
	if version != "" {
		common.SetVersion(version)
	}

	cmd.Execute()
}
