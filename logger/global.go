package logger

import "io"

var logger *Logger

func Panic(format string, args ...interface{}) {
	logger.Panic(format, args...)
}

func Fatal(format string, args ...interface{}) {
	logger.Fatal(format, args...)
}

func Error(format string, args ...interface{}) {
	logger.Error(format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.Warn(format, args...)
}

func Info(format string, args ...interface{}) {
	logger.Info(format, args...)
}

func Debug(format string, args ...interface{}) {
	logger.Debug(format, args...)
}

func Trace(format string, args ...interface{}) {
	logger.Trace(format, args...)
}

func SetName(name string) {
	logger.SetName(name)
}

func SetLevel(level Level) {
	logger.SetLevel(level)
}

func SetOutput(output io.Writer) {
	logger.SetOutput(output)
}

func SetFormat(format Format) {
	logger.SetFormat(format)
}

func AddHook(hook Hook) {
	logger.AddHook(hook)
}

func init() {
	logger = New()
}
