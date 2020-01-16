package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// type Logger struct {
// 	*zerolog.Logger
// }

// func New(isErr, isInfo, isDebug bool) *Logger {
// 	logLevel := zerolog.InfoLevel
// 	if isDebug {
// 		logLevel = zerolog.DebugLevel
// 	} else if isInfo {
// 		logLevel = zerolog.InfoLevel
// 	} else if isErr {
// 		logLevel = zerolog.ErrorLevel
// 	} else {
// 		logLevel = zerolog.FatalLevel
// 	}

// 	zerolog.SetGlobalLevel(logLevel)
// 	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
// 	return &Logger{&logger}
// }

// func NewConsole(isErr, isInfo, isDebug bool) *Logger {
// 	logLevel := zerolog.InfoLevel
// 	if isDebug {
// 		logLevel = zerolog.DebugLevel
// 	} else if isInfo {
// 		logLevel = zerolog.InfoLevel
// 	} else if isErr {
// 		logLevel = zerolog.ErrorLevel
// 	} else {
// 		logLevel = zerolog.FatalLevel
// 	}

// 	zerolog.SetGlobalLevel(logLevel)
// 	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
// 	return &Logger{&logger}
// }

func New(isErr, isInfo, isDebug bool, callerInfo bool) *zerolog.Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	} else if isInfo {
		logLevel = zerolog.InfoLevel
	} else if isErr {
		logLevel = zerolog.ErrorLevel
	} else {
		logLevel = zerolog.FatalLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	var logger zerolog.Logger
	if callerInfo {
		logger = zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	} else {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
	return &logger
}

func NewConsole(isErr, isInfo, isDebug bool, callerInfo bool) *zerolog.Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	} else if isInfo {
		logLevel = zerolog.InfoLevel
	} else if isErr {
		logLevel = zerolog.ErrorLevel
	} else {
		logLevel = zerolog.FatalLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	var logger zerolog.Logger
	if callerInfo {
		logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	return &logger
}

// // Output duplicates the global logger and sets w as its output.
// func (l *Logger) Output(w io.Writer) zerolog.Logger {
// 	return l.Output(w)
// }
