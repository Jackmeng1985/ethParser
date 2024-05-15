package log

import (
	"fmt"
	"time"
)

// LogLevel defines the severity of the log message
type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// Logger represents a simple logger
type Logger struct {
	level LogLevel
}

// New creates a new logger instance
func New(level LogLevel) *Logger {
	return &Logger{
		level: level,
	}
}

// log writes a log message with the given level and message
func (l *Logger) log(level LogLevel, msg string) {
	if level >= l.level {
		timestamp := time.Now().Format(time.RFC3339)
		levelStr := ""
		color := ""
		switch level {
		case INFO:
			levelStr = "INFO"
			color = colorBlue
		case WARN:
			levelStr = "WARN"
			color = colorYellow
		case ERROR:
			levelStr = "ERROR"
			color = colorRed
		}

		//n := runtime.Callers(4, pcs[:])

		logMsg := fmt.Sprintf("%s  %s[%s]%s %s\n", timestamp, color, levelStr, colorReset, msg)
		fmt.Print(logMsg) // Print to console
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, ctx ...interface{}) {
	l.log(INFO, fmt.Sprintf(msg, ctx...))
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, ctx ...interface{}) {
	l.log(WARN, fmt.Sprintf(msg, ctx...))
}

// Error logs an error message
func (l *Logger) Error(msg string, ctx ...interface{}) {
	l.log(ERROR, fmt.Sprintf(msg, ctx...))
}
