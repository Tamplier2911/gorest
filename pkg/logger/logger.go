package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New is used to create a new logger instance.
func New(logLevel string, production bool) *zap.SugaredLogger {
	// create .og level
	var level zapcore.Level
	level.Set(logLevel)

	// logger config
	config := zap.Config{
		Development:      !production,
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeDuration: zapcore.SecondsDurationEncoder,
			LevelKey:       "severity",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			CallerKey:      "caller",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			TimeKey:        "timestamp",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			NameKey:        "name",
			EncodeName:     zapcore.FullNameEncoder,
			MessageKey:     "message",
			StacktraceKey:  "",
			LineEnding:     "\n",
		},
	}

	if !production {
		config.Encoding = "console"
		// config.EncoderConfig.LevelKey = ""
		// config.EncoderConfig.CallerKey = ""
		config.EncoderConfig.TimeKey = ""
		config.EncoderConfig.NameKey = ""
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// build logger from config
	logger, _ := config.Build()

	// configure and create logger
	return logger.Sugar()
}
