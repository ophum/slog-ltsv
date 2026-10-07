package main

import (
	"errors"
	"log/slog"
	"os"

	slogltsv "github.com/ophum/slog-ltsv"
)

func main() {
	handler := slogltsv.NewLTSVHandler(os.Stdout, slogltsv.Option{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler).With(
		slog.String("service", "checkout"),
		slog.Int("version", 1),
	)

	logger.Info("server started", slog.String("addr", ":8080"))

	requestLogger := logger.WithGroup("request").With(
		slog.String("id", "req-123"),
		slog.String("method", "GET"),
	)
	requestLogger.Debug("request received", slog.String("path", "/orders"))

	err := errors.New("upstream timeout")
	requestLogger.Error("request failed",
		slog.String("path", "/orders"),
		slog.Any("error", err),
	)

	logger.Warn("input contains LTSV delimiters",
		slog.String("input", "first\tsecond\nthird\\"),
	)
}
