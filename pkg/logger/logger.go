package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Level represents the severity level of a log message
type Level string

const (
	// Debug level for development information
	Debug Level = "DEBUG"
	// Info level for general operational information
	Info Level = "INFO"
	// Warn level for warnings
	Warn Level = "WARN"
	// Error level for errors
	Error Level = "ERROR"
	// Fatal level for critical errors
	Fatal Level = "FATAL"
)

// Logger defines a structured logger interface
type Logger interface {
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(msg string, err error, fields ...map[string]interface{})
	Fatal(msg string, err error, fields ...map[string]interface{})
}

// JSONLogger implements Logger using JSON format
type JSONLogger struct {
	serviceName string
}

// NewJSONLogger creates a new JSON logger
func NewJSONLogger(serviceName string) Logger {
	return &JSONLogger{
		serviceName: serviceName,
	}
}

// log outputs a structured log message
func (l *JSONLogger) log(level Level, msg string, err error, fields ...map[string]interface{}) {
	logEntry := map[string]interface{}{
		"timestamp":    time.Now().Format(time.RFC3339),
		"level":        level,
		"service_name": l.serviceName,
		"message":      msg,
	}

	// Add error if present
	if err != nil {
		logEntry["error"] = err.Error()
	}

	// Add additional fields
	if len(fields) > 0 {
		for key, value := range fields[0] {
			logEntry[key] = value
		}
	}

	// Marshal to JSON
	jsonLog, jsonErr := json.Marshal(logEntry)
	if jsonErr != nil {
		fmt.Fprintf(os.Stderr, "error marshaling log entry: %v\n", jsonErr)
		return
	}

	// Write to stdout for regular logs, stderr for errors and fatal
	if level == Error || level == Fatal {
		fmt.Fprintln(os.Stderr, string(jsonLog))
	} else {
		fmt.Fprintln(os.Stdout, string(jsonLog))
	}

	// Exit on fatal errors
	if level == Fatal {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *JSONLogger) Debug(msg string, fields ...map[string]interface{}) {
	l.log(Debug, msg, nil, fields...)
}

// Info logs an informational message
func (l *JSONLogger) Info(msg string, fields ...map[string]interface{}) {
	l.log(Info, msg, nil, fields...)
}

// Warn logs a warning message
func (l *JSONLogger) Warn(msg string, fields ...map[string]interface{}) {
	l.log(Warn, msg, nil, fields...)
}

// Error logs an error message
func (l *JSONLogger) Error(msg string, err error, fields ...map[string]interface{}) {
	l.log(Error, msg, err, fields...)
}

// Fatal logs a fatal error message and exits
func (l *JSONLogger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	l.log(Fatal, msg, err, fields...)
}
