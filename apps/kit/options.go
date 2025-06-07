package kit

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/welltop-cn/common/cloud/tracer"
	"github.com/welltop-cn/common/utils/log"
)

type kitOptions struct {
	serviceName  string
	shutdownFunc []func(ctx context.Context) error
}

func (k *kitOptions) waitingShutdown() {
	defer func() {
		if err := recover(); err != nil {
			log.CtxErrorf(context.Background(), "panic error, %v", err)
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	log.CtxInfof(context.Background(), "receive signal, start to shutdown")
	for index, f := range k.shutdownFunc {
		log.CtxInfof(context.Background(), "shutdownFunc index: %d", index)
		err := f(context.Background())
		if err != nil {
			log.CtxErrorf(context.Background(), "shutdown error, %v", err)
		}
	}
}

type KitOptions func(o *kitOptions)

func WithTracer() KitOptions {
	return func(o *kitOptions) {
		log.CtxInfof(context.Background(), "==== init otel tracing ===")
		shutdown, err := tracer.InitOtelTracer()
		if err != nil {
			log.CtxErrorf(context.Background(), "failed to init otel tracer ", err)
		}
		if shutdown != nil {
			if len(o.shutdownFunc) <= 0 {
				o.shutdownFunc = make([]func(ctx context.Context) error, 0, 5)
			}
			o.shutdownFunc = append(o.shutdownFunc, shutdown)
			log.CtxInfof(context.Background(), "=== init otel tracing success ===")
		}
	}
}
