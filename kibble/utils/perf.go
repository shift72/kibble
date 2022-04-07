//    Copyright 2018 SHIFT72
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

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
func NewStopwatchfWithLevel(msg string, a ...interface{}) *Stopwatch {
	return &Stopwatch{msg: fmt.Sprintf(msg, a...), start: time.Now(), level: logging.NOTICE}
}

// NewStopwatchf - start a stop watch with formatting
func NewStopwatchf(msg string, a ...interface{}) *Stopwatch {
	return &Stopwatch{msg: fmt.Sprintf(msg, a...), start: time.Now(), level: logging.DEBUG}
}

// Completed - stops the stop watch
func (sw *Stopwatch) Completed() time.Duration {
	d := round(time.Since(sw.start), time.Millisecond)
	if log.IsEnabledFor(sw.level) {
		log.Noticef("%s: %s", sw.msg, d)
	} else {
		log.Debugf("%s: %s", sw.msg, d)
	}
	return d
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
