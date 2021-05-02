package logger

import (
	"os"

	"github.com/rs/zerolog"
	"golang.org/x/term"
)

// Level defines logger level.
type Level int8

// Available log levels.
const (
	Debug = iota
	Info
)

// Logger writes logs.
type Logger struct {
	logger zerolog.Logger
}

// New returns new logger.
func New(level Level) *Logger {
	outFile := os.Stdout

	logger := zerolog.New(outFile)

	if term.IsTerminal(int(outFile.Fd())) {
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out: outFile,
		})
	}

	logger = logger.With().
		Timestamp().
		CallerWithSkipFrameCount(3).
		Logger().
		Level(zerolog.Level(level))

	return &Logger{
		logger: logger,
	}
}

// Debug writes debug log.
func (l *Logger) Debug(format string, vv ...interface{}) {
	l.logger.Debug().Msgf(format, vv...)
}

// Info writes information log.
func (l *Logger) Info(format string, vv ...interface{}) {
	l.logger.Info().Msgf(format, vv...)
}

// Error writes error log.
func (l *Logger) Error(format string, vv ...interface{}) {
	l.logger.Error().Msgf(format, vv...)
}

// Fatal writes fatal log and stops the application.
func (l *Logger) Fatal(format string, vv ...interface{}) {
	l.logger.Fatal().Msgf(format, vv...)
}
