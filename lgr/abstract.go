package lgr

// Logger is interface to be implemented
// on different layers of application.
type Logger interface {
	Debugf(string, ...any)
	Debugw(string, ...any)
	Infof(string, ...any)
	Infow(string, ...any)
	Warnf(string, ...any)
	Warnw(string, ...any)
	Errorf(string, ...any)
	Errorw(string, ...any)
}
