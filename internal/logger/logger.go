package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init() {

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	Log = zerolog.New(output).
		With().
		Timestamp().
		Logger()
}

func Info() *zerolog.Event {
	return Log.Info()
}

func Warn() *zerolog.Event {
	return Log.Warn()
}

func Error(err error) *zerolog.Event {

	if err == nil {
		return Log.Error()
	}

	return Log.Error().
		Err(err)
}

func Fatal(err error) *zerolog.Event {

	if err == nil {
		return Log.Fatal()
	}

	return Log.Fatal().
		Err(err)
}
