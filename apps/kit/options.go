package kit

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type kitOptions struct {
	serviceName  string
	shutdownFunc []func(ctx context.Context) error
}

func (k *kitOptions) waitingShutdown() {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("panic error, %v", err)
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	slog.Info("receive signal, start to shutdown")
	for index, f := range k.shutdownFunc {
		slog.Info("shutdownFunc index: %d", index)
		err := f(context.Background())
		if err != nil {
			slog.Error("shutdown error，%v", err)
		}
	}
}

type KitOptions func(o *kitOptions)

func WithTracer() KitOptions {
	return func(o *kitOptions) {
		slog.Info("==== init otel tracing ===")
		shutdown, err := tracer.InitOtelTracer()
		if err != nil {
			slog.Error("failed to init otel tracer ", err)
		}
		if shutdown != nil {
			if len(o.shutdownFunc) <= 0 {
				o.shutdownFunc = make([]func(ctx context.Context) error, 0, 5)
			}
			o.shutdownFunc = append(o.shutdownFunc, shutdown)
			slog.Info("=== init otel tracing success ===")
		}
	}
}
