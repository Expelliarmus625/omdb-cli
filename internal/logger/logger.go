package logger

import (
	"log/slog"
	"os"
	"path/filepath"
)

var Log *slog.Logger

func Init(logPath string) error {
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return err
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}

	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	Log = slog.New(handler)
	slog.SetDefault(Log)

	return nil
}
