package logger

import (
    "log/slog"
    "os"
    "strings"
)

func New(level string) *slog.Logger {
    lvl := slog.LevelInfo
    switch strings.ToUpper(level) {
    case "DEBUG": lvl = slog.LevelDebug
    case "WARN": lvl = slog.LevelWarn
    case "ERROR": lvl = slog.LevelError
    }
    return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
}
