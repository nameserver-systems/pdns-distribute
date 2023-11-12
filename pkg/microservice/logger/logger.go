package logger

import (
	"os"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var signalChannel chan<- os.Signal //nolint:gochecknoglobals

// InitGlobalLogger set global logging settings.
func InitGlobalLogger(sigchannel chan<- os.Signal) {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	signalChannel = sigchannel
}

// SetDefaultLogLevel with panic, fatal, error, warn, info, debug
// default: info
func SetDefaultLogLevel(level string) {
	var loglevel zerolog.Level

	switch level {
	case "panic":
		loglevel = zerolog.PanicLevel
	case "fatal":
		loglevel = zerolog.FatalLevel
	case "error":
		loglevel = zerolog.ErrorLevel
	case "warn":
		loglevel = zerolog.WarnLevel
	case "info":
		loglevel = zerolog.InfoLevel
	case "debug":
		loglevel = zerolog.DebugLevel
	default:
		loglevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(loglevel)
}

// Log panic, fatal, error, warn, info, debug
// default: info
func Log(level, message string, v ...interface{}) {
	switch level {
	case "panic":
		log.Panic().Msgf(message, v...)
	case "fatal":
		log.Fatal().Msgf(message, v...)
	case "error":
		log.Error().Msgf(message, v...)
	case "warn":
		log.Warn().Msgf(message, v...)
	case "info":
		log.Info().Msgf(message, v...)
	case "debug":
		log.Debug().Msgf(message, v...)
	default:
		log.Info().Msgf(message, v...)
	}
}

func DebugLog(message string) {
	log.Debug().Msg(message)
}

func InfoLog(message string) {
	log.Info().Msg(message)
}

func WarnLog(message string) {
	log.Warn().Msg(message)
}

func ErrorLog(message string) {
	log.Error().Msg(message)
}

func FatalLog(message string) {
	signalChannel <- syscall.SIGINT

	log.Log().Str("level", "fatal").Msg(message)
}

func PanicLog(message string) {
	signalChannel <- syscall.SIGINT

	log.Panic().Msg(message)
}

func DebugErrLog(err error) {
	log.Debug().Err(err).Msg("")
}

func InfoErrLog(err error) {
	log.Info().Err(err).Msg("")
}

func WarnErrLog(err error) {
	log.Warn().Err(err).Msg("")
}

func ErrorErrLog(err error) {
	log.Error().Err(err).Msg("")
}

func FatalErrLog(err error) {
	signalChannel <- syscall.SIGINT

	log.Log().Str("level", "fatal").Err(err).Msg("")
}

func PanicErrLog(err error) {
	signalChannel <- syscall.SIGINT

	log.Panic().Err(err).Msg("")
}
