package ent

// Logger is designed for logging.
type Logger interface {
	Flush() error
	Level() string
	Debugf(string, ...any)
	Debugw(string, ...any)
	Infof(string, ...any)
	Infow(string, ...any)
	Warnf(string, ...any)
	Warnw(string, ...any)
	Errorf(string, ...any)
	Errorw(string, ...any)
	// Fatal ... are essential for migrations.
	Fatal(...any)
	Fatalf(string, ...any)
	Print(...any)
	Println(...any)
	Printf(string, ...any)
}