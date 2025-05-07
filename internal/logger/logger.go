package logger

import (
	"fmt"
	"log/slog"
	"os"
)

// InitLogger создаёт и настраивает глобальный slog.Logger.
// level — минимальный уровень, format — "text" или "json".
func InitLogger(level slog.Level, format string) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		AddSource: true,  // печатает файл:строку, откуда лог вызван
		Level:     level, // минимальный уровень логов
	}

	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default: //something wrong. regardless of the check in the config
		fmt.Fprintf(os.Stderr, "invalid LOG_FORMAT %q, must be one of [text, json]\n", format)
		os.Exit(1)
	}

	logger := slog.New(handler)
	// делаем его глобальным (slog.Default())
	slog.SetDefault(logger)
	return logger
}
