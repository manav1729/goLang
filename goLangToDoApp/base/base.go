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

var ctx context.Context

func Init() {
	// Set default logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx = context.WithValue(context.Background(), TraceIDString, uuid.New())
}

func Exit() {
	// Signal Channel listens for
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	// System Exit when signal received
	go func() {
		<-signalChannel
		LogInfo("Received termination signal, shutting down...")
		os.Exit(0)
	}()

	// Infinite loop to keep the application running
	LogInfo("Application is Running. Press Ctrl+C to exit.")
	select {}
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
