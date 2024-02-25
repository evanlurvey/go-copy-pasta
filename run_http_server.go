package main

import (
	"context"
	"net"
	"net/http"
	"time"
)

// creates a server that will shut down when the context is cancelled or a failure occurs
// this function will block until the server is shut down
func RunServer(ctx context.Context, srv http.Handler) {
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("localhost", "8080"),
		Handler: srv,
	}
	failure := make(chan error)
	go func() {
		// l.Info("http server starting", zap.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			failure <- err
		}
	}()
	wait := make(chan struct{})
	go func() {
		defer close(wait)
		select {
		case <-ctx.Done():
			// l.Info("http server shutting down")
			shutdownCtx := context.Background()
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if err := httpServer.Shutdown(shutdownCtx); err != nil {
				_ = err
				// l.Error("http server shutdown failed", zap.Error(err))
			}
		case err := <-failure:
			_ = err
			// l.Error("http server start failed", zap.Error(err))
		}
	}()
	<-wait
}
