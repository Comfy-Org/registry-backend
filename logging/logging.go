package drip_logging

import (
	"context"
	"os"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func SetupLogger() zerolog.Logger {
	var once sync.Once
	var log zerolog.Logger

	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		// Needed to conform with GCP Cloud Logging format.
		zerolog.LevelFieldName = "severity"
		zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
			switch l {
			case zerolog.DebugLevel:
				return "DEBUG"
			case zerolog.InfoLevel:
				return "INFO"
			case zerolog.WarnLevel:
				return "WARNING"
			case zerolog.ErrorLevel:
				return "ERROR"
			case zerolog.FatalLevel:
				return "CRITICAL"
			case zerolog.PanicLevel:
				return "ALERT"
			default:
				return "DEFAULT"
			}
		}
		log = zerolog.New(os.Stdout)

		if os.Getenv("DRIP_ENV") == "localdev" {
			log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}
	})

	return log
}

// ReuseContextLogger returns a new context with the same logger as the given context
func ReuseContextLogger(ctx context.Context, newCtx context.Context) context.Context {
	l := zerolog.Ctx(ctx)
	if l == nil {
		return newCtx
	}
	return l.WithContext(newCtx)
}

func SetGlobalLogLevel(logLevel string) {
	// Default to info level
	defaultLevel := zerolog.InfoLevel
	if logLevel != "" {
		level, err := zerolog.ParseLevel(logLevel)
		if err == nil {
			defaultLevel = level
		} else {
			log.Error().Err(err).Msg("Invalid log level")
		}
	}

	zerolog.SetGlobalLevel(defaultLevel)
}
