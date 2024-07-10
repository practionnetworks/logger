// config.go

package logger

import (
	"errors"
)

type Config struct {
	ServiceName        string
	Pod                string
	LogLevel           string // Log level as string (e.g., "Debug", "Info", etc.)
	LogAnalyserAddress string // Optional, set to nil if not used
	LogAnalyserEnabled bool   // Optional, set to true if not used
	Console            bool   // Optional, set to false if not used
	LogFilePath        string // Optional, leave empty if not used
}

func NewLogger(serviceName string, console bool, pod string, logFilePath string, logAnalyserAddress string, logLevel string, LogAnalyserEnabled bool) (Config, error) {
	if !console && logFilePath == "" && logAnalyserAddress == "" {
		return Config{}, errors.New("at least one logging option (Console, LogFile, LogAnalyserAddress) must be selected")
	}

	if serviceName == "" && pod == "" {
		return Config{}, errors.New("service name and PodName must be provided")
	}

	return Config{
		ServiceName:        serviceName,
		Pod:                pod,
		LogLevel:           logLevel,
		LogAnalyserAddress: logAnalyserAddress,
		LogAnalyserEnabled: LogAnalyserEnabled,
		Console:            console,
		LogFilePath:        logFilePath,
	}, nil
}
