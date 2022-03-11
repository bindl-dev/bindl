package internal

import (
	"github.com/rs/zerolog"
	"go.xargs.dev/bindl/internal/log"
)

func Log() *zerolog.Logger {
	return log.Log
}
