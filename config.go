// config.go

package logger

import (
	"errors"
	"io"
)

type Config struct {
	ServiceName string
	Pod         string
	LogLevel    string    // Log level as string (e.g., "Debug", "Info", etc.)
	LogAnalyser io.Writer // Optional, set to nil if not used
	Console     bool      // Optional, set to false if not used
	LogFilePath string    // Optional, leave empty if not used
}

func DefaultConfig(serviceName string, console bool, pod string, logFilePath string, logAnalyser io.Writer, logLevel string) (Config, error) {
	if !console && logFilePath == "" && logAnalyser == nil {
		return Config{}, errors.New("at least one logging option (Console, LogFile, LogAnalyser) must be selected")
	}

	if serviceName == "" && pod == "" {
		return Config{}, errors.New("service name and PodName must be provided")
	}

	return Config{
		ServiceName: serviceName,
		Pod:         pod,
		LogLevel:    logLevel,
		Console:     console,
		LogFilePath: logFilePath,
		LogAnalyser: logAnalyser,
	}, nil
}
