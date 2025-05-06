package logger

import (
	"github.com/rs/zerolog/log"
)

var Logger = log.Logger

func init() {
}

func Error(err error) {
	Logger.Err(err).Send()
}
