package log

import (
	"go.uber.org/zap/zapcore"
)

type Level string

const (
	Debug Level = "debug"
	Info  Level = "info"
	WARN  Level = "warn"
	Error Level = "error"
	Fatal Level = "fatal"
)

var levelZapMap = map[Level]zapcore.Level{
	Debug: zapcore.DebugLevel,
	Info:  zapcore.InfoLevel,
	WARN:  zapcore.WarnLevel,
	Error: zapcore.ErrorLevel,
	Fatal: zapcore.FatalLevel,
}

func Level2ZapLevle(l Level) zapcore.Level {
	if zapl, ok := levelZapMap[l]; ok {
		return zapl
	}
	return zapcore.InfoLevel
}
