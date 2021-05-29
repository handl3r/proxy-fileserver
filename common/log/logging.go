package log

type Logging interface {
	Info(args ...interface{})
	Infof(msg string, args ...interface{})
	Warn(args ...interface{})
	Warnf(msg string, args ...interface{})
	Error(args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(msg string, args ...interface{})
}
