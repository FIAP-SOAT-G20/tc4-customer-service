package logger

import (
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/infrastructure/config"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(cfg *config.Config) *Logger {
	var handler slog.Handler

	if cfg.Environment == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})
	} else {
		handler = NewPrettyHandler(os.Stdout, PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			}})
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}
