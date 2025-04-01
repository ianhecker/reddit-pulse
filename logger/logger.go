package logger

import "log"

type Logger struct {
	Verbose bool
}

func MakeLogger() Logger {
	return Logger{true}
}

func (l *Logger) SetVerbose(verbose bool) {
	l.Verbose = verbose
}

func (l *Logger) Log(format string, a ...any) {
	if !l.Verbose {
		return
	}
	log.Printf(format, a...)
}
