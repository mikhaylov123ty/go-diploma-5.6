package logger

import (
	"context"
	"github.com/google/uuid"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
	"log/slog"
	"os"
	"strings"
)

type customLogger struct {
	slog.Handler
}

func Init(level string) {
	slogLevel := parseLevel(strings.ToLower(level))
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel, AddSource: false})
	pretty := &customLogger{handler}

	slog.SetDefault(slog.New(pretty))
}

func (c *customLogger) Handle(ctx context.Context, rec slog.Record) error {
	if v := ctx.Value(utils.ContextKey("login")); v != nil {
		rec.Add("login", v.(string))
	}
	if v := ctx.Value(utils.ContextKey("processID")); v != nil {
		rec.Add("processID", v.(uuid.UUID).String())
	}
	_ = c.Handler.Handle(ctx, rec)
	return nil
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo

	}
}
