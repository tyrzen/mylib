package lgr

import (
	"log"

	"go.uber.org/zap"
)

type Level string

const (
	Debug Level = "DEBUG"
	Info  Level = "INFO"
	Warn  Level = "WARN"
	Error Level = "ERROR"
)

type Logger struct{ *zap.Logger }

func New(text Level) Logger {
	enc := zap.NewDevelopmentEncoderConfig()

	lvl, err := zap.ParseAtomicLevel(string(text))
	if err != nil {
		log.Fatalf("failed logger level: %w", err)
	}

	cfg := zap.Config{
		Level:            lvl,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    enc,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("failed creating logger: %w", err)
	}

	return Logger{logger}
}
