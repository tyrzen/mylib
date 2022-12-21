package lgr

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level string

const (
	Debug Level = "DEBUG"
	Info  Level = "INFO"
	Warn  Level = "WARN"
	Error Level = "ERROR"
)

type Log struct{ core *zap.SugaredLogger }

func New(level Level) Log {
	lvl, err := zap.ParseAtomicLevel(string(level))
	if err != nil {
		log.Fatalf("failed logger level: %v", err)
	}

	enc := zapcore.EncoderConfig{
		MessageKey:          "message",
		LevelKey:            "level",
		TimeKey:             "timestamp",
		NameKey:             "name",
		CallerKey:           "caller",
		FunctionKey:         "",
		StacktraceKey:       "stacktrace",
		SkipLineEnding:      false,
		LineEnding:          "\n",
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          zapcore.ISO8601TimeEncoder,
		EncodeDuration:      zapcore.StringDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "",
	}

	cfg := zap.Config{
		Level:            lvl,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    enc,
		OutputPaths:      []string{"stdout", "/app.log"},
		ErrorOutputPaths: []string{"stderr", "/app.log"},
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("failed creating logger: %v", err)
	}

	return Log{logger.Sugar()}
}

func (log *Log) Debugf(format string, keyVal ...any) {
	log.core.Debugf(format, keyVal...)
}

func (log *Log) Debugw(format string, keyVal ...any) {
	log.core.Debugw(format, keyVal...)
}

func (log *Log) Infof(format string, keyVal ...any) {
	log.core.Infof(format, keyVal...)
}

func (log *Log) Infow(format string, keyVal ...any) {
	log.core.Infow(format, keyVal...)
}

func (log *Log) Warnf(format string, keyVal ...any) {
	log.core.Warnf(format, keyVal...)
}

func (log *Log) Warnw(format string, keyVal ...any) {
	log.core.Warnw(format, keyVal...)
}

func (log *Log) Errorf(format string, keyVal ...any) {
	log.core.Errorf(format, keyVal...)
}

func (log *Log) Errorw(format string, keyVal ...any) {
	log.core.Errorw(format, keyVal...)
}
