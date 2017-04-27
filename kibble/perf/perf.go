package perf

import (
	"fmt"
	"time"
)

// Stopwatch - track a stop watch
type Stopwatch struct {
	msg   string
	start time.Time
}

// NewStopwatch - start a stop watch
func NewStopwatch(msg string) Stopwatch {
	return Stopwatch{msg: msg, start: time.Now()}
}

// NewStopwatchf - start a stop watch with formatting
func NewStopwatchf(msg string, a ...interface{}) Stopwatch {
	return Stopwatch{msg: fmt.Sprintf(msg, a), start: time.Now()}
}

// Completed - stops the stop watch
func (sw *Stopwatch) Completed() {
	fmt.Printf("%s: %s\n", sw.msg, time.Now().Sub(sw.start))
}

// MeasureElapsed - measure the time taken to complete the function
func MeasureElapsed(msg string, fn func()) {
	start := time.Now()
	fn()
	stop := time.Now()
	fmt.Printf("%s: %s\n", msg, stop.Sub(start).String())
}
