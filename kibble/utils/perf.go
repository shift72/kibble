package utils

import (
	"fmt"
	"time"

	logging "github.com/op/go-logging"
)

// Stopwatch - track a stop watch
type Stopwatch struct {
	msg   string
	level logging.Level
	start time.Time
}

// NewStopwatch - start a stop watch
func NewStopwatch(msg string) *Stopwatch {
	return &Stopwatch{msg: msg, start: time.Now(), level: logging.DEBUG}
}

// NewStopwatchLevel - start a stop watch with a level
func NewStopwatchLevel(msg string, level logging.Level) *Stopwatch {
	return &Stopwatch{msg: msg, start: time.Now(), level: level}
}

// NewStopwatchf - start a stop watch with formatting
func NewStopwatchf(msg string, a ...interface{}) *Stopwatch {
	return &Stopwatch{msg: fmt.Sprintf(msg, a), start: time.Now(), level: logging.DEBUG}
}

// Completed - stops the stop watch
func (sw *Stopwatch) Completed() {
	if log.IsEnabledFor(sw.level) {
		log.Noticef("%s: %s", sw.msg, round(time.Now().Sub(sw.start), time.Millisecond))
	} else {
		log.Debugf("%s: %s", sw.msg, round(time.Now().Sub(sw.start), time.Millisecond))
	}
}

// MeasureElapsed - measure the time taken to complete the function
func MeasureElapsed(msg string, fn func()) {
	start := time.Now()
	fn()
	stop := time.Now()
	if log.IsEnabledFor(logging.NOTICE) {
		log.Noticef("%s: %s", msg, stop.Sub(start).String())
	} else {
		log.Debugf("%s: %s", msg, stop.Sub(start).String())
	}
}

func round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}
