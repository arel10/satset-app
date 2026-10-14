package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(env string) *Logger {
	var zlog zerolog.Logger

	if env == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		zlog = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		zlog = zerolog.New(output).With().Timestamp().Logger()
	}

	return &Logger{Logger: zlog}
}
