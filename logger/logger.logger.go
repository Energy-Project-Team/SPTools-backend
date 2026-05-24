// Copyright (C) 2026 Energy Project Team
// SPDX-License-Identifier: AGPL-3.0-only
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

var GlobalLogger *slog.Logger

type LoggerConfig struct {
	AddSource bool
	NoColor   bool
	Level     slog.Level
	LogDir    string
}

func InitLogger(cfg LoggerConfig) {
	if cfg.LogDir == "" {
		cfg.LogDir = "logs"
	}

	writers := []io.Writer{os.Stdout}

	if err := os.MkdirAll(cfg.LogDir, 0o755); err == nil {
		if file, err := os.OpenFile(filepath.Join(cfg.LogDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
			writers = append(writers, file)
		}
		if file, err := os.OpenFile(filepath.Join(cfg.LogDir, "combined.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
			writers = append(writers, file)
		}
	}

	handler := tint.NewHandler(io.MultiWriter(writers...), &tint.Options{
		AddSource:  cfg.AddSource,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    cfg.NoColor,
		Level:      cfg.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.SourceKey {
				return a
			}

			src, ok := a.Value.Any().(*slog.Source)
			if !ok || src == nil {
				return a
			}

			file := filepath.Base(src.File)
			line := strconv.Itoa(src.Line)

			value := file + ":" + line

			if src.Function != "" {
				fn := shortFuncName(src.Function)
				value += " (" + fn + ")"
			}

			return slog.String(slog.SourceKey, value)
		},
	})

	GlobalLogger = slog.New(handler)
}

func logWithCaller(level slog.Level, msg string, args ...any) {
	if GlobalLogger == nil {
		InitLogger(LoggerConfig{
			AddSource: true,
			NoColor:   false,
			Level:     slog.LevelDebug,
		})
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	record := slog.NewRecord(time.Now(), level, msg, pcs[0])
	record.Add(args...)

	_ = GlobalLogger.Handler().Handle(context.Background(), record)
}

func Info(msg string, args ...any) {
	logWithCaller(slog.LevelInfo, msg, args...)
}

func Debug(msg string, args ...any) {
	logWithCaller(slog.LevelDebug, msg, args...)
}

func Warn(msg string, args ...any) {
	logWithCaller(slog.LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	logWithCaller(slog.LevelError, msg, args...)
}

func Fatal(msg string, args ...any) {
	logWithCaller(slog.LevelError, msg, args...)
	os.Exit(1)
}

func shortFuncName(full string) string {
	full = strings.ReplaceAll(full, "\\", "/")

	if idx := strings.LastIndex(full, "/"); idx >= 0 {
		full = full[idx+1:]
	}

	parts := strings.Split(full, ".")
	if len(parts) <= 1 {
		return full
	}

	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}

	return full
}
