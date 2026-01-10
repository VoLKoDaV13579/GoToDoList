package logger

import (
	"log/slog"
	"os"
)

// Logger представляет обертку над структурированным логгером
type Logger struct {
	logger *slog.Logger
}

// New создает новый экземпляр логгера
func New(debug bool) *Logger {
	var level slog.Level

	// Устанавливаем уровень логирования в зависимости от режима
	if debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	// Создаем handler для JSON-логирования (идеально для production)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		// Добавляем информацию о месте вызова лога (файл и строка)
		AddSource: debug,
	})

	return &Logger{
		logger: slog.New(handler),
	}
}

// NewTextLogger создает логгер с текстовым выводом (удобнее для разработки)
func NewTextLogger(debug bool) *Logger {
	var level slog.Level

	if debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: debug,
	})

	return &Logger{
		logger: slog.New(handler),
	}
}

// Debug логирует сообщение уровня DEBUG
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

// Info логирует сообщение уровня INFO
func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

// Warn логирует сообщение уровня WARNING
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

// Error логирует сообщение уровня ERROR
func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

// With создает новый логгер с дополнительными полями контекста
func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		logger: l.logger.With(args...),
	}
}
