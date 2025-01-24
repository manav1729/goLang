package base

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"os"
)

const TraceIDString = "trace_id"

var ctx context.Context

func Init() {
	// Set default logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx = context.WithValue(context.Background(), TraceIDString, uuid.New())
}

func LogInfo(msg string, args ...any) {
	args = appendArgs(args)
	slog.InfoContext(ctx, msg, args...)
}

func LogWarn(msg string, args ...any) {
	args = appendArgs(args)
	slog.WarnContext(ctx, msg, args...)
}

func LogDebug(msg string, args ...any) {
	args = appendArgs(args)
	slog.DebugContext(ctx, msg, args...)
}

func LogError(msg string, args ...any) {
	args = appendArgs(args)
	slog.ErrorContext(ctx, msg, args...)
}

func appendArgs(args []any) []any {
	return append(args, TraceIDString, getTraceID())
}

func getTraceID() uuid.UUID {
	return ctx.Value(TraceIDString).(uuid.UUID)
}
