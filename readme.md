# CLI Timer
CLI tool that is very similar to `sleep` but it shows a timer with the remaining time (timer is removed when done).

## Install
You can install the package using [go](https://go.dev/doc/install) with:
```bash
go install github.com/AmrSaber/timer
```

If you want to uninstall it, you can run
```bash
rm $(which timer)
```

## Usage
Use as:
```
timer 10s
```

You can add multiple durations (including negative ones) and they will all be summed together to get the final timer duration.

Durations can have any of go duration suffexis (h, m, s, ms, us, ...) but whatever duration is given it is rounded to the closest second. If no suffix is given to a duration, then `s` is assumed.

Not providing any duration, providing wrong duration, or cancelling the timer before it's done (sending SIGINT signal with ctrl+c) results in exit code of 1.

