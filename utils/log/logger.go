package log

import (
	"context"
	"fmt"
	"os"

	"github.com/welltop-cn/common/cloud/config"
	"github.com/welltop-cn/common/protos"
	"github.com/welltop-cn/common/utils/env"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type SpanContext struct {
	TraceId      trace.TraceID
	SpanId       trace.SpanID
	ParentSpanId trace.SpanID
	SessionId    string
	OneId        string
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
	enc.AddString("session_id", t.SessionId)
	enc.AddString("one_id", t.OneId)
	enc.AddString("path", t.Path)
	return nil
}

func getLoggerConf() *protos.Logger {
	var logC protos.Logger
	err := config.Get("logger").Scan(&logC)
	if err != nil {
		panic("get logger config failed" + err.Error())
	}
	return &logC
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.MessageKey = "body"
	encoderConfig.EncodeName = zapcore.FullNameEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.NameKey = ""
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 使用 lumberjack 库设置log归档、切分
func getLumberJackLogger(fileLog *protos.Logger) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", fileLog.Path, fileLog.Filename),
		MaxSize:    0,
		MaxAge:     0,
		MaxBackups: 0,
		Compress:   true,
		LocalTime:  true,
	}
}

func InitLogger() {
	logConf := getLoggerConf()
	encoder := getEncoder()

	lumberJackLogger := getLumberJackLogger(logConf)

	multiWriter := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(Level2ZapLevle(Level(logConf.Level)))

	core := zapcore.NewCore(
		encoder,     // 设置编码器
		multiWriter, // 设置日志打印方式
		atomicLevel, // 日志级别
	)

	caller := zap.AddCaller()                                     // 开启开发模式，堆栈跟踪
	development := zap.Development()                              // 开启文件及行号
	filed := zap.Fields(zap.String("service", env.ServiceName())) // 设置初始化字段
	skip := zap.AddCallerSkip(1)
	logger := zap.New(core, caller, development, filed, skip)

	zap.ReplaceGlobals(logger)
}

func buildCtxField(ctx context.Context) zap.Field {
	// span := oteltrace.SpanFromContext(ctx)
	// if !span.IsRecording() {
	// 	return emptyField
	// }
	// if rspan, ok := span.(sdktrace.ReadOnlySpan); ok {
	// 	traceId := rspan.SpanContext().TraceID()
	// 	spanId := rspan.SpanContext().SpanID()
	// 	var pid oteltrace.SpanID
	// 	if rspan.Parent().IsValid() {
	// 		pid = rspan.Parent().SpanID()
	// 	}

	// 	var tangoCtx = &SpanContext{
	// 		TraceId:      traceId,
	// 		SpanId:       spanId,
	// 		ParentSpanId: pid,
	// 		Path:         rspan.Name(),
	// 	}

	// 	return zap.Object(CtxField, tangoCtx)
	// }

	return zap.Object(CtxField, &SpanContext{
		SessionId: "sss",
	})
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
	zap.L().Debug(msg, buildCtxField(ctx))
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
