package log

import (
	"context"
	"fmt"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SpanContext struct {
	TraceId      trace.TraceID
	SpanId       trace.SpanID
	ParentSpanId trace.SpanID
	Path         string
}

func (t *SpanContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if t.TraceId.IsValid() {
		enc.AddString("trace_id", t.TraceId.String())
	}
	if t.SpanId.IsValid() {
		enc.AddString("span_id", t.SpanId.String())
	}
	if t.ParentSpanId.IsValid() {
		enc.AddString("parent_span_id", t.ParentSpanId.String())
	}
	enc.AddString("path", t.Path)
	return nil
}

func buildCtxField(ctx context.Context) zap.Field {
	if ctx == nil {
		return zap.Object(CtxField, &SpanContext{})
	}
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return zap.Object(CtxField, &SpanContext{})
	}
	if rspan, ok := span.(sdktrace.ReadOnlySpan); ok {
		traceId := rspan.SpanContext().TraceID()
		spanId := rspan.SpanContext().SpanID()
		var pid trace.SpanID
		if rspan.Parent().IsValid() {
			pid = rspan.Parent().SpanID()
		}

		var tangoCtx = &SpanContext{
			TraceId:      traceId,
			SpanId:       spanId,
			ParentSpanId: pid,
			Path:         rspan.Name(),
		}

		return zap.Object(CtxField, tangoCtx)
	}

	return zap.Object(CtxField, &SpanContext{})
}

func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

func CtxDebug(ctx context.Context, a ...interface{}) {
	msg := getMessage("", a)
	zap.L().Log(zap.DebugLevel, msg, buildCtxField(ctx))
}

func CtxDebugf(ctx context.Context, format string, a ...interface{}) {
	msg := getMessage(format, a)
	zap.L().Debug(msg, buildCtxField(ctx))

}

// func CtxDebugw(ctx context.Context, msg string, keyvalues ...interface{}) {
// 	if log.innerLogger.shouldLog() {
// 		fields := log.sweetenFields(keyvalues)
// 		log.innerLogger.zlog.Debug(msg, log.innerLogger.addTangoCtx(ctx, fields)...)
// 	}
// }

func CtxInfo(ctx context.Context, a ...interface{}) {
	msg := getMessage("", a)
	zap.L().Info(msg, buildCtxField(ctx))
}

// // CtxInfof uses fmt.Sprintf to log a templated message.
func CtxInfof(ctx context.Context, format string, a ...interface{}) {
	msg := getMessage(format, a)
	zap.L().Info(msg, buildCtxField(ctx))
}

// func CtxInfow(ctx context.Context, msg string, keyvalues ...interface{}) {
// 	if log.innerLogger.shouldLog() {
// 		fields := log.sweetenFields(keyvalues)
// 		log.innerLogger.zlog.Info(msg, log.innerLogger.addTangoCtx(ctx, fields)...)
// 	}
// }

func CtxWarn(ctx context.Context, a ...interface{}) {
	msg := getMessage("", a)
	zap.L().Warn(msg, buildCtxField(ctx))
}

func CtxWarnf(ctx context.Context, format string, a ...interface{}) {
	msg := getMessage(format, a)
	zap.L().Warn(msg, buildCtxField(ctx))
}

// func CtxWarnw(ctx context.Context, msg string, keyvalues ...interface{}) {
// 	if log.innerLogger.shouldLog() {
// 		fields := log.sweetenFields(keyvalues)
// 		log.innerLogger.zlog.Warn(msg, log.innerLogger.addTangoCtx(ctx, fields)...)
// 	}
// }

func CtxError(ctx context.Context, a ...interface{}) {
	msg := getMessage("", a)
	zap.L().Error(msg, buildCtxField(ctx))
}

func CtxErrorf(ctx context.Context, format string, a ...interface{}) {
	msg := getMessage(format, a)
	zap.L().Error(msg, buildCtxField(ctx))
}

// func CtxErrorw(ctx context.Context, msg string, keyvalues ...interface{}) {
// 	if log.innerLogger.shouldLog() {
// 		fields := log.sweetenFields(keyvalues)
// 		log.innerLogger.zlog.Error(msg, log.innerLogger.addTangoCtx(ctx, fields)...)
// 	}
// }

func CtxFatal(ctx context.Context, a ...interface{}) {
	msg := getMessage("", a)
	zap.L().Fatal(msg, buildCtxField(ctx))
}

func CtxFatalf(ctx context.Context, format string, a ...interface{}) {
	msg := getMessage(format, a)
	zap.L().Fatal(msg, buildCtxField(ctx))
}

// func CtxFatalw(ctx context.Context, msg string, keyvalues ...interface{}) {
// 	if log.innerLogger.shouldLog() {
// 		fields := log.sweetenFields(keyvalues)
// 		log.innerLogger.zlog.Fatal(msg, log.innerLogger.addTangoCtx(ctx, fields)...)
// 	}
// }
