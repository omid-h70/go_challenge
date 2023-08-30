package worker

import (
	"fmt"
	"github.com/rs/zerolog"
)

type Logger struct {
}

// we implemented asynq.Logger interface
func NewLogger() *Logger {
	return &Logger{}
}

func (l *logger) Print(int logLevel, args ...any) {
	log.WithLevel(logLevel).Msg(fmt.Sprint(args...))
}

func (l *logger) Debug(args ...any) {
	l.Print(zerolog.DebugLevel, args...)
}

func (l *logger) Info(args ...any) {
	l.Print(zerolog.InfoLevel, args...)
}

func (l *logger) Warn(args ...any) {
	l.Print(zerolog.WarnLevel, args...)
}

func (l *logger) Error(args ...any) {
	l.Print(zerolog.ErrorLevel, args...)
}

func (l *logger) Fatal(args ...any) {
	l.Print(zerolog.FatalLevel, args...)
}
