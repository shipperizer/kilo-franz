package testkit

import (
	"fmt"
)

// FakeLogger is a simple in-memory logger used in tests.
type FakeLogger struct {
	LastMessage string
	LastMsg     string
	Message     string
	Messages    []string
	LastLevel   string
}

func (l *FakeLogger) update(level, message string) {
	if l == nil {
		return
	}
	l.LastLevel = level
	l.LastMessage = message
	l.LastMsg = message
	l.Message = message
	l.Messages = append(l.Messages, message)
}

func (l *FakeLogger) Errorf(format string, args ...interface{}) {
	l.update("error", fmt.Sprintf(format, args...))
}

func (l *FakeLogger) Infof(format string, args ...interface{}) {
	l.update("info", fmt.Sprintf(format, args...))
}

func (l *FakeLogger) Warnf(format string, args ...interface{}) {
	l.update("warn", fmt.Sprintf(format, args...))
}

func (l *FakeLogger) Debugf(format string, args ...interface{}) {
	l.update("debug", fmt.Sprintf(format, args...))
}

func (l *FakeLogger) Fatalf(format string, args ...interface{}) {
	l.update("fatal", fmt.Sprintf(format, args...))
}

func (l *FakeLogger) Error(args ...interface{}) {
	l.update("error", fmt.Sprint(args...))
}

func (l *FakeLogger) Info(args ...interface{}) {
	l.update("info", fmt.Sprint(args...))
}

func (l *FakeLogger) Warn(args ...interface{}) {
	l.update("warn", fmt.Sprint(args...))
}

func (l *FakeLogger) Debug(args ...interface{}) {
	l.update("debug", fmt.Sprint(args...))
}

func (l *FakeLogger) Fatal(args ...interface{}) {
	l.update("fatal", fmt.Sprint(args...))
}
