package base

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const TraceIDString = "trace_id"

func Init() context.Context {
	// Set default logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.WithValue(context.Background(), TraceIDString, uuid.New())
	return ctx
}

func Exit(ctx context.Context) {
	// Signal Channel listens for
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)

	// System Exit when signal received
	go func() {
		<-signalChannel
		LogInfo(ctx, "Received termination signal, shutting down...")
		os.Exit(0)
	}()

	// Infinite loop to keep the application running
	LogInfo(ctx, "Application is Running. Press Ctrl+C to exit.")
	select {}
}

func LogInfo(ctx context.Context, msg string, args ...any) {
	args = appendArgs(ctx, args)
	slog.InfoContext(ctx, msg, args...)
}

func LogDebug(ctx context.Context, msg string, args ...any) {
	args = appendArgs(ctx, args)
	slog.DebugContext(ctx, msg, args...)
}

func LogError(ctx context.Context, msg string, args ...any) {
	args = appendArgs(ctx, args)
	slog.ErrorContext(ctx, msg, args...)
}

func appendArgs(ctx context.Context, args []any) []any {
	return append(args, TraceIDString, getTraceID(ctx))
}

func getTraceID(ctx context.Context) uuid.UUID {
	return ctx.Value(TraceIDString).(uuid.UUID)
}
