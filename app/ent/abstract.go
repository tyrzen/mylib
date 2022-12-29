package ent

// Logger is interface to be implemented
// on different layers of application.
type Logger interface {
	Flush() error
	Level() string
	Debugf(string, ...interface{})
	Debugw(string, ...interface{})
	Infof(string, ...interface{})
	Infow(string, ...interface{})
	Warnf(string, ...interface{})
	Warnw(string, ...interface{})
	Errorf(string, ...interface{})
	Errorw(string, ...interface{})
	// Fatal ... are essential for migrations.
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Print(...interface{})
	Println(...interface{})
	Printf(string, ...interface{})
}
