package logger

import (
	"github.com/sirupsen/logrus"
	"io"
)

type Level = logrus.Level

const (
	PanicLevel Level = logrus.PanicLevel
	FatalLevel Level = logrus.FatalLevel
	ErrorLevel Level = logrus.ErrorLevel
	WarnLevel  Level = logrus.WarnLevel
	InfoLevel  Level = logrus.InfoLevel
	DebugLevel Level = logrus.DebugLevel
	TraceLevel Level = logrus.TraceLevel
)

const defName = "fun_log"

type Logger struct {
	logger *logrus.Logger
	name   string
}

func New(opts ...Option) *Logger {
	log := &Logger{
		logger: logrus.New(),
		name:   defName,
	}

	log.SetName(defName)
	log.SetLevel(InfoLevel)
	log.SetOutput(Stdout())
	log.SetFormat(FunFormat(nil))

	for _, opt := range opts {
		opt(log)
	}

	return log
}

func (l *Logger) SetName(name string) {
	l.name = name
}

func (l *Logger) SetLevel(level Level) {
	l.logger.SetLevel(level)
}

func (l *Logger) SetOutput(output io.Writer) {
	l.logger.SetOutput(output)
}

func (l *Logger) SetFormat(format Format) {
	l.logger.SetFormatter(format)
}

func (l *Logger) AddHook(hook Hook) {
	l.logger.AddHook(hook)
}

func (l *Logger) Panic(format string, args ...interface{}) {
	l.logger.WithField("name", l.name).Panicf(format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.logger.WithField("name", l.name).Fatalf(format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.logger.WithField("name", l.name).Errorf(format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.logger.WithField("name", l.name).Warnf(format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.logger.WithField("name", l.name).Infof(format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.logger.WithField("name", l.name).Debugf(format, args...)
}

func (l *Logger) Trace(format string, args ...interface{}) {
	l.logger.WithField("name", l.name).Tracef(format, args...)
}
