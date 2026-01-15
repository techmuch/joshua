package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Init initializes the global slog logger
func Init(path, level string) {
	var programLevel = slog.LevelInfo
	switch strings.ToUpper(level) {
	case "DEBUG":
		programLevel = slog.LevelDebug
	case "INFO":
		programLevel = slog.LevelInfo
	case "WARN":
		programLevel = slog.LevelWarn
	case "ERROR":
		programLevel = slog.LevelError
	}

	// File writer with rotation
	fileWriter := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	// MultiWriter to log to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)

	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: programLevel,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
