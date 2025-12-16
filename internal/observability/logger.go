package observability

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger for consistent logging across the application
type Logger struct {
	logger zerolog.Logger
}

// NewLogger creates a new logger instance with the specified log level
func NewLogger(level string) *Logger {
	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel // default to info
	}

	// Configure zerolog
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano // ISO 8601 with nanoseconds

	// Use JSON format for production (structured logging)
	// Use ConsoleWriter for development (human-readable)
	useConsole := os.Getenv("LOG_FORMAT") == "console"

	var output zerolog.ConsoleWriter
	if useConsole {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	logger := zerolog.New(output).With().
		Timestamp().
		Logger()

	return &Logger{logger: logger}
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Infof logs an info message with fields
func (l *Logger) Infof(msg string, fields map[string]interface{}) {
	event := l.logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Errorf logs an error message with fields and optional error
func (l *Logger) Errorf(msg string, err error, fields map[string]interface{}) {
	event := l.logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Warnf logs a warning message with fields
func (l *Logger) Warnf(msg string, fields map[string]interface{}) {
	event := l.logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf logs a debug message with fields
func (l *Logger) Debugf(msg string, fields map[string]interface{}) {
	event := l.logger.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// WithRequestID creates a logger instance with a request ID field
// This is crucial for tracing requests across distributed systems
func (l *Logger) WithRequestID(requestID string) *Logger {
	newLogger := l.logger.With().Str("request_id", requestID).Logger()
	return &Logger{logger: newLogger}
}

// WithField creates a logger instance with additional fields
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := l.logger.With().Interface(key, value).Logger()
	return &Logger{logger: newLogger}
}

// GetZerologLogger returns the underlying zerolog.Logger for advanced usage
func (l *Logger) GetZerologLogger() zerolog.Logger {
	return l.logger
}

