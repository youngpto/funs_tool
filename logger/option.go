package logger

import "io"

type Option func(logger *Logger)

func WithName(name string) Option {
	return func(logger *Logger) {
		logger.SetName(name)
	}
}

func WithLevel(level Level) Option {
	return func(logger *Logger) {
		logger.SetLevel(level)
	}
}

func WithOutput(output io.Writer) Option {
	return func(logger *Logger) {
		logger.SetOutput(output)
	}
}

func WithFormat(format Format) Option {
	return func(logger *Logger) {
		logger.SetFormat(format)
	}
}

func WithHook(hook Hook) Option {
	return func(logger *Logger) {
		logger.AddHook(hook)
	}
}
