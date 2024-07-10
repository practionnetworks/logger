// logging.go

package logger

import (
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var initialized bool

type LogstashWriter struct {
	conn net.Conn
}

func NewLogstashWriter(network, address string) (*LogstashWriter, error) {
	conn, err := net.Dial(network, address)

	if err != nil {
		return nil, err
	}
	return &LogstashWriter{conn: conn}, nil
}

func (w *LogstashWriter) Write(p []byte) (n int, err error) {
	return w.conn.Write(p)
}

func InitLogger(config Config) {
	if initialized {
		log.Warn().Msg("Logger already initialized, skipping re-initialization")
		return
	}

	zerolog.TimeFieldFormat = time.RFC3339

	var writers []io.Writer

	// Add console output if enabled
	if config.Console {
		// writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}) // Disable ANSI escape codes

		writers = append(writers, os.Stdout)
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

	if config.LogAnalyserEnabled {
		logstashWriter, err := NewLogstashWriter("tcp", config.LogAnalyserAddress)

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create Logstash writer")
		}

		writers = append(writers, logstashWriter)
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
	logLevel := parseLogLevel(config.LogLevel)

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

func parseLogLevel(level string) zerolog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return zerolog.DebugLevel
	case "INFO":
		return zerolog.InfoLevel
	case "WARN":
		return zerolog.WarnLevel
	case "ERROR":
		return zerolog.ErrorLevel
	case "FATAL":
		return zerolog.FatalLevel
	case "PANIC":
		return zerolog.PanicLevel
	case "TRACE":
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}
func logWithFields(level zerolog.Level, message string, fields ...interface{}) {
	event := log.WithLevel(level)
	if len(fields)%2 != 0 {
		event = event.Interface("fields_error", "uneven number of key-value pairs")
	} else {
		for i := 0; i < len(fields); i += 2 {
			key, okKey := fields[i].(string)
			value, okValue := fields[i+1].(string)
			if okKey && okValue {
				event = event.Str(key, value)
			} else {
				event = event.Interface("fields_error", "key-value pairs must be strings")
				break
			}
		}
	}
	event.Msg(message)
}

func Info(message string, fields ...interface{}) {
	logWithFields(zerolog.InfoLevel, message, fields...)
}

func Debug(message string, fields ...interface{}) {
	logWithFields(zerolog.DebugLevel, message, fields...)
}

func Warn(message string, fields ...interface{}) {
	logWithFields(zerolog.WarnLevel, message, fields...)
}

func Error(message string, fields ...interface{}) {
	logWithFields(zerolog.ErrorLevel, message, fields...)
}

func Fatal(message string, fields ...interface{}) {
	logWithFields(zerolog.FatalLevel, message, fields...)
}

func Panic(message string, fields ...interface{}) {
	logWithFields(zerolog.PanicLevel, message, fields...)
}

func Trace(message string, fields ...interface{}) {
	logWithFields(zerolog.TraceLevel, message, fields...)
}

func WarnWithError(err error, fields ...interface{}) {
	logWithFields(zerolog.WarnLevel, err.Error(), append(fields, "error", errors.WithStack(err))...)
}

func ErrorWithError(err error, fields ...interface{}) {
	logWithFields(zerolog.ErrorLevel, err.Error(), append(fields, "error", errors.WithStack(err))...)
}

func FatalWithError(err error, fields ...interface{}) {
	logWithFields(zerolog.FatalLevel, err.Error(), append(fields, "error", errors.WithStack(err))...)
}

func PanicWithError(err error, fields ...interface{}) {
	logWithFields(zerolog.PanicLevel, err.Error(), append(fields, "error", errors.WithStack(err))...)
}

func TraceWithError(err error, fields ...interface{}) {
	logWithFields(zerolog.TraceLevel, err.Error(), append(fields, "error", errors.WithStack(err))...)
}

// func Info(message string, fields ...Field) {
// 	event := log.Info()
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(message)
// }

// func Debug(message string, fields ...Field) {
// 	event := log.Debug()
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(message)
// }

// func Warn(message string, fields ...Field) {
// 	event := log.Warn()
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(message)
// }

// func Error(message string, fields ...Field) {
// 	event := log.Error()
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(message)
// }

// func Fatal(message string, fields ...Field) {
// 	event := log.Fatal()
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(message)
// }

// func Panic(message string, fields ...Field) {
// 	event := log.Panic()
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(message)
// }

// func Trace(message string, fields ...Field) {
// 	event := log.Trace()
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(message)
// }

// func WarnWithError(err error, fields ...Field) {
// 	event := log.Warn().Stack().Err(errors.WithStack(err))
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(err.Error())
// }

// func ErrorWithError(err error, fields ...Field) {
// 	event := log.Error().Stack().Err(errors.WithStack(err))
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(err.Error())
// }

// func FatalWithError(err error, fields ...Field) {
// 	event := log.Fatal().Stack().Err(errors.WithStack(err))
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(err.Error())
// }

// func PanicWithError(err error, fields ...Field) {
// 	event := log.Panic().Stack().Err(errors.WithStack(err))
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(err.Error())
// }

// func TraceWithError(err error, fields ...Field) {
// 	event := log.Trace().Stack().Err(errors.WithStack(err))
// 	for _, field := range fields {
// 		event = event.Str(field.Key, field.Value)
// 	}
// 	event.Msg(err.Error())
// }

// // Field represents a key-value pair for structured logging
// type Field struct {
// 	Key   string
// 	Value string
// }
