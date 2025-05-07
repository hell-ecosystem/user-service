package logger

import (
	"log/slog"
	"os"
)

// InitLogger создаёт и настраивает глобальный slog.Logger.
// Его стоит вызывать в самом начале приложения (в cmd/serve.go и cmd/migrate.go).
func InitLogger(level slog.Level) *slog.Logger {
	// TextHandler — человекочитаемый вывод, можно заменить на NewJSONHandler
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,  // печатает файл:строку, откуда лог вызван
		Level:     level, // минимальный уровень логов
	})
	logger := slog.New(handler)
	// делаем его глобальным (slog.Default())
	slog.SetDefault(logger)
	return logger
}
