package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Use console writer for development
	if os.Getenv("ENV") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	// Add service name to all logs
	log.Logger = log.With().Str("service", "test-wallet").Logger()
}

// LogError logs an error with context
func LogError(err error, msg string, fields map[string]interface{}) {
	log.Error().Err(err).Fields(fields).Msg(msg)
}

// LogInfo logs an info message with context
func LogInfo(msg string, fields map[string]interface{}) {
	log.Info().Fields(fields).Msg(msg)
}

// LogDebug logs a debug message with context
func LogDebug(msg string, fields map[string]interface{}) {
	log.Debug().Fields(fields).Msg(msg)
}

// LogFatal logs a fatal error and exits
func LogFatal(err error, msg string, fields map[string]interface{}) {
	log.Fatal().Err(err).Fields(fields).Msg(msg)
}
