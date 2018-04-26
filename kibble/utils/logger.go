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
func ConfigureStandardLogging(verbose bool) {
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
	setLogLevel(verbose)
}

// ConfigureInteractiveLogging - logging
func ConfigureInteractiveLogging(verbose bool) {
	logging.SetFormatter(
		logging.MustStringFormatter(
			`%{color}%{message}%{color:reset}`,
		))
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
	setLogLevel(verbose)
}

// ConfigureWatchedLogging - logging to stdout + the unique logger
func ConfigureWatchedLogging(verbose bool) *UniqueLogger {
	uni := NewUniqueLogger()
	log1 := logging.NewBackendFormatter(uni, logging.MustStringFormatter(
		`%{level} - %{message}`,
	))
	log2 := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetBackend(logging.MultiLogger(log1, log2))
	setLogLevel(verbose)
	return uni
}

// ConfigureSyncLogging - logging only to the unique logger
func ConfigureSyncLogging(verbose bool) *UniqueLogger {
	uni := NewUniqueLogger()
	logging.SetBackend(
		logging.NewBackendFormatter(uni,
			logging.MustStringFormatter(`%{level} - %{message}`),
		),
	)
	setLogLevel(verbose)
	return uni
}

func setLogLevel(verbose bool) {
	if verbose {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}
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
func NewUniqueLogger() *UniqueLogger {
	return &UniqueLogger{
		level: logging.WARNING,
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
