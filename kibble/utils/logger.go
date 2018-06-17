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
	"os"

	logging "github.com/op/go-logging"
)

func init() {
	logging.SetFormatter(
		logging.MustStringFormatter(
			`%{color}%{time:15:04:05.000} ▶ %{message}%{color:reset}`,
		))
}

// ConfigureStandardLoggingLevel - logging
func ConfigureStandardLoggingLevel(level logging.Level) {
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
	logging.SetLevel(level, "")
}

// ConfigureStandardLogging - verbose
func ConfigureStandardLogging(level logging.Level) {
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
	logging.SetLevel(level, "")
}

// ConfigureInteractiveLogging - logging
func ConfigureInteractiveLogging(level logging.Level) {
	logging.SetFormatter(
		logging.MustStringFormatter(
			`%{color}%{message}%{color:reset}`,
		))
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
	logging.SetLevel(level, "")
}

// ConfigureWatchedLogging - logging to stdout + the unique logger
func ConfigureWatchedLogging(level logging.Level) *UniqueLogger {
	uni := NewUniqueLogger(logging.WARNING)
	log1 := logging.NewBackendFormatter(uni, logging.MustStringFormatter(
		`%{level} - %{message}`,
	))
	log2 := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetBackend(logging.MultiLogger(log1, log2))
	logging.SetLevel(level, "")
	return uni
}

// ConfigureSyncLogging - logging only to the unique logger
func ConfigureSyncLogging(level logging.Level) *UniqueLogger {
	uni := NewUniqueLogger(level)
	logging.SetBackend(
		logging.NewBackendFormatter(uni,
			logging.MustStringFormatter(`%{level} - %{message}`),
		),
	)
	logging.SetLevel(level, "")
	return uni
}

// ConvertToLoggingLevel - convert to
func ConvertToLoggingLevel(verbose bool) logging.Level {
	if verbose {
		return logging.DEBUG
	}
	return logging.INFO
}

// UniqueLogger - logs only the unique errors
type UniqueLogger struct {
	level logging.Level
	store []string
}

// LogReader - reads logs
type LogReader interface {
	Logs() []string
	Clear()
}

// NewUniqueLogger - creates a new unique logger
func NewUniqueLogger(level logging.Level) *UniqueLogger {
	return &UniqueLogger{
		level: level,
		store: make([]string, 0),
	}
}

// Log - only the unique errors
func (l *UniqueLogger) Log(level logging.Level, calldepth int, rec *logging.Record) (err error) {

	if level > l.level {
		return
	}

	nm := rec.Formatted(calldepth)
	for _, m := range l.store {
		if nm == m {
			return
		}
	}

	l.store = append(l.store, nm)
	return
}

// Logs -
func (l *UniqueLogger) Logs() []string {
	return l.store
}

// Clear the logs
func (l *UniqueLogger) Clear() {
	l.store = make([]string, 0)
}
