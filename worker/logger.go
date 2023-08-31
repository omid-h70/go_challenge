package worker

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
}

// NewLogger we implemented asynq.Logger interface
func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Print(logLevel zerolog.Level, args ...any) {
	log.WithLevel(logLevel).Msg(fmt.Sprint(args...))
}

func (l *Logger) Debug(args ...any) {
	l.Print(zerolog.DebugLevel, args...)
}

func (l *Logger) Info(args ...any) {
	l.Print(zerolog.InfoLevel, args...)
}

func (l *Logger) Warn(args ...any) {
	l.Print(zerolog.WarnLevel, args...)
}

func (l *Logger) Error(args ...any) {
	l.Print(zerolog.ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...any) {
	l.Print(zerolog.FatalLevel, args...)
}
