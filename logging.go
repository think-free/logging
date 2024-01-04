/*
	This is a simple logging helper that use zap and permit
	- to add tags to the log
	- to save the logger in the context to allow the tags to be added in the next log
	(C)2023 - Christophe Meurice (meumeu1402@gmail.com)
*/

package logging

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type t string

var lgstr t = "logger"
var defaultLogger *zap.Logger

func init() {
	Init("debug")
}

func Init(l string) {
	level := getZapLevelFromString(l)

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	defaultLogger, err = config.Build()
	if err != nil {
		panic(err)
	}
	// Set StdLog and globals to use the default logger
	_ = zap.RedirectStdLog(defaultLogger)
	_ = zap.ReplaceGlobals(defaultLogger)

	defaultLogger.Sugar().Infof("logging initialized with level '%s'", l)
}

type Lg struct {
	*zap.SugaredLogger
}

func ContextWithLogger(ctx context.Context) context.Context {
	if ctx.Value(lgstr) != nil {
		return ctx // Already have a logger, don't overwrite it
	}
	logger := Logger(ctx)
	return context.WithValue(ctx, lgstr, logger)
}

func Logger(ctx context.Context) *Lg {
	// If the context already have a logger, return it
	// Else return the default logger
	if ctx.Value(lgstr) != nil {
		return ctx.Value(lgstr).(*Lg)
	}
	return &Lg{
		defaultLogger.Sugar(),
	}
}

func (l *Lg) SetTag(key string, value interface{}) {
	l.SugaredLogger = l.With(key, value)
}

func getZapLevelFromString(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	case "panic":
		return zap.PanicLevel
	}
	return zap.InfoLevel
}

// Use this function if you don't want the key/value pair to be present in the context logger 
// Use it if you set the same tag repeatedly in a loop for example as the tag will be duplicated each time
func GetZapLoggerWithValue(ctx context.Context, key string, value interface{}) *zap.SugaredLogger {
	return Logger(ctx).With(key, value)
}
