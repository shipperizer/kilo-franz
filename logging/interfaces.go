package logging

type LoggerInterface interface {
	Errorf(string, ...any)
	Infof(string, ...any)
	Warnf(string, ...any)
	Debugf(string, ...any)
	Fatalf(string, ...any)
	Error(...any)
	Info(...any)
	Warn(...any)
	Debug(...any)
	Fatal(...any)
}
