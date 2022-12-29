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

type Log struct{ core *zap.SugaredLogger }

func New() *Log {
	lvl := os.Getenv("LOG_LEVEL")
	if lvl == "" {
		lvl = DebugLevel
	}

	atomLvl, err := zap.ParseAtomicLevel(lvl)
	if err != nil {
		log.Panicf("failed to parse logger level: %v", err)
	}

	encConsoleCfg := zapcore.EncoderConfig{
		MessageKey:          "message",
		LevelKey:            "level",
		TimeKey:             "time",
		NameKey:             "name",
		CallerKey:           "caller",
		FunctionKey:         "",
		StacktraceKey:       "stacktrace",
		SkipLineEnding:      false,
		LineEnding:          "\n",
		EncodeLevel:         zapcore.CapitalColorLevelEncoder,
		EncodeTime:          zapcore.ISO8601TimeEncoder,
		EncodeDuration:      zapcore.NanosDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "\t",
	}

	encFileCfg := encConsoleCfg
	encFileCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	file, err := os.Create(os.Getenv("LOG_FILE"))
	if err != nil {
		log.Fatalf("failed creating logger file : %v", err)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(zapcore.NewConsoleEncoder(encConsoleCfg), zapcore.Lock(os.Stderr), atomLvl.Level()),
		zapcore.NewCore(zapcore.NewJSONEncoder(encFileCfg), zapcore.Lock(file), atomLvl.Level()),
	}

	core := zapcore.NewTee(cores...)

	return &Log{zap.New(core).Sugar()}
}

func (log *Log) Level() string {
	return log.core.Level().String()
}

func (log *Log) Flush() error {
	return log.core.Sync()
}

func (log *Log) Debugf(format string, keyVal ...interface{}) {
	log.core.Debugf(format, keyVal...)
}

func (log *Log) Debugw(format string, keyVal ...interface{}) {
	log.core.Debugw(format, keyVal...)
}

func (log *Log) Infof(format string, keyVal ...interface{}) {
	log.core.Infof(format, keyVal...)
}

func (log *Log) Infow(format string, keyVal ...interface{}) {
	log.core.Infow(format, keyVal...)
}

func (log *Log) Warnf(format string, keyVal ...interface{}) {
	log.core.Warnf(format, keyVal...)
}

func (log *Log) Warnw(format string, keyVal ...interface{}) {
	log.core.Warnw(format, keyVal...)
}

func (log *Log) Errorf(format string, keyVal ...interface{}) {
	log.core.Errorf(format, keyVal...)
}

func (log *Log) Errorw(format string, keyVal ...interface{}) {
	log.core.Errorw(format, keyVal...)
}

func (log *Log) Fatal(keyVal ...interface{}) {
	log.core.Fatal(keyVal...)
}

func (log *Log) Fatalf(format string, keyVal ...interface{}) {
	log.core.Fatalf(format, keyVal...)
}

func (log *Log) Print(keyVal ...interface{}) {
	log.core.Info(keyVal...)
}

func (log *Log) Println(keyVal ...interface{}) {
	log.core.Infoln(keyVal...)
}

func (log *Log) Printf(format string, keyVal ...interface{}) {
	log.core.Infof(format, keyVal...)
}
