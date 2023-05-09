package utils

import (
	"github.com/rs/zerolog"
	"os"
)

var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

func GetLogger() *zerolog.Logger {
	return &logger
}
