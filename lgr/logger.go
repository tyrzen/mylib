package lgr

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "DEBUG"
	InfoLevel  = "INFO"
	WarnLevel  = "WARN"
	ErrorLevel = "ERROR"
)

const FilePath = "app.log"

type Log struct{ core *zap.SugaredLogger }

func New(level string) *Log {
	file, err := os.Create(FilePath)
	if err != nil {
		log.Fatalf("failed creating logger file : %v", err)
	}

	defer file.Close()

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		log.Panicf("failed logger level: %v", err)
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
		EncodeLevel:         zapcore.CapitalColorLevelEncoder,
		EncodeTime:          zapcore.ISO8601TimeEncoder,
		EncodeDuration:      zapcore.StringDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "\n",
	}

	cfg := zap.Config{
		Level:            lvl,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    enc,
		OutputPaths:      []string{"stdout", FilePath},
		ErrorOutputPaths: []string{"stderr", FilePath},
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("failed build logger: %v", err)
	}

	return &Log{logger.Sugar()}
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
