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
const DataFile = "../data/ToDoData.json"

type customHandler struct {
	slog.Handler
}

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID, ok := ctx.Value(TraceIDString).(string); ok {
		r.AddAttrs(slog.String(TraceIDString, traceID))
	}
	if traceID, ok := ctx.Value(TraceIDString).(uuid.UUID); ok {
		r.AddAttrs(slog.String(TraceIDString, traceID.String()))
	}
	return h.Handler.Handle(ctx, r)
}

func Init() context.Context {
	// Set default logger
	baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler := &customHandler{baseHandler}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	ctx := context.WithValue(context.Background(), TraceIDString, uuid.New())
	return ctx
}

func Exit(ctx context.Context) {
	// Signal channel listens for
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)

	// System Exit go routine
	go systemExit(ctx, signalChannel)

	// Infinite loop to keep the application running
	slog.InfoContext(ctx, "Application is Running. Press Ctrl+C to exit.")
	select {}
}

// Basic system exit go routine
func systemExit(ctx context.Context, signalChannel chan os.Signal) {
	<-signalChannel
	slog.InfoContext(ctx, "Received termination signal, shutting down...")
	os.Exit(0)
}
