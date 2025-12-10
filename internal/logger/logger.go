package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

// Logger wraps the charmbracelet log logger.
type Logger struct {
	*log.Logger
}

// New creates a new logger instance.
func New(verbose bool) *Logger {
	level := log.InfoLevel
	if verbose {
		level = log.DebugLevel
	}

	l := log.New(os.Stderr)
	l.SetLevel(level)
	return &Logger{Logger: l}
}

// SetLevel sets the log level.
func (l *Logger) SetLevel(level log.Level) {
	l.Logger.SetLevel(level)
}

