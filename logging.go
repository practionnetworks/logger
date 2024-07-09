// logging.go

package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var initialized bool

func InitLogger(config Config) {
	if initialized {
		log.Warn().Msg("Logger already initialized, skipping re-initialization")
		return
	}

	zerolog.TimeFieldFormat = time.RFC3339

	var writers []io.Writer

	// Add console output if enabled
	if config.Console {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}) // Disable ANSI escape codes
	}

	// Add file output if provided
	if config.LogFilePath != "" {
		file, err := os.OpenFile(config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open log file")
		}
		writers = append(writers, file)

		// Store file handle in a package-level variable to ensure it's not closed prematurely
		log.Logger = log.Logger.Output(file)
	}

	// Add log analyzer output if provided
	if config.LogAnalyser != nil {
		writers = append(writers, config.LogAnalyser)
	}

	// Use MultiWriter to combine outputs
	var multiWriter io.Writer
	if len(writers) > 0 {
		multiWriter = io.MultiWriter(writers...)
	} else {
		// Default to stdout if no specific output configured
		multiWriter = os.Stdout
	}

	// Convert log level string to zerolog.Level
	var logLevel zerolog.Level
	switch strings.ToLower(config.LogLevel) {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	case "fatal":
		logLevel = zerolog.FatalLevel
	case "panic":
		logLevel = zerolog.PanicLevel
	case "trace":
		logLevel = zerolog.TraceLevel
	default:
		// Default to Info if no valid log level is provided
		logLevel = zerolog.InfoLevel
	}

	// Initialize logger with JSON formatter
	log.Logger = zerolog.New(multiWriter).With().
		Timestamp().
		Str("service", config.ServiceName).
		Str("pod", config.Pod).
		Int("pid", os.Getpid()).
		CallerWithSkipFrameCount(3).
		Logger().
		Level(logLevel).
		Output(multiWriter) // Use multiWriter for output

	initialized = true
}

func Info(message string) {
	log.Info().Msg(message)
}

func Debug(message string) {
	log.Debug().Msg(message)
}

func Warn(message string) {
	log.Warn().Msg(message)
}

func Error(message string) {
	log.Error().Msg(message)
}

func Fatal(message string) {
	log.Fatal().Msg(message)
}

func Panic(message string) {
	log.Panic().Msg(message)
}

func Trace(message string) {
	log.Trace().Msg(message)
}

func WarnWithError(err error) {
	log.Warn().Stack().Err(errors.WithStack(err)).Msg(err.Error())
}

func ErrorWithError(err error) {
	log.Error().Stack().Err(errors.WithStack(err)).Msg(err.Error())
}

func FatalWithError(err error) {
	log.Fatal().Stack().Err(errors.WithStack(err)).Msg(err.Error())
}

func PanicWithError(err error) {
	log.Panic().Stack().Err(errors.WithStack(err)).Msg(err.Error())
}

func TraceWithError(err error) {
	log.Trace().Stack().Err(errors.WithStack(err)).Msg(err.Error())
}
