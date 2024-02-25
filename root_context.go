package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// creates a root context that will start a shutdown when a signal is received
func RootContext() context.Context {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// l.Info("signal received, shutting down...")
		cancel()
	}()
	// ctx = logging.WithLogger(ctx, l)
	return ctx
}
