package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var _ FullLogger = (*DefaultLogger)(nil)

const (
	LogIDKey = "request_id"
	ReqIDKey = "X-Request-ID"
)

// DefaultLogger is a wrapper around `zerolog.Logger` that provides an implementation of `FullLogger` interface
type DefaultLogger struct {
	log     zerolog.Logger
	out     io.Writer
	level   zerolog.Level
	options []Opt
}

// ConsoleWriter parses the JSON input and writes it in an
// (optionally) colorized, human-friendly format to Out.
type ConsoleWriter = zerolog.ConsoleWriter

func NewConsole(Out io.Writer) ConsoleWriter {
	cw := ConsoleWriter{}
	cw.NoColor = true
	cw.Out = Out
	cw.TimeFormat = "2006-01-02 15:04:05"
	cw.FormatLevel = func(lv interface{}) string {
		return fmt.Sprintf("[%s]", lv)
	}
	// 自定义 Caller 显示格式，只显示文件名和行数
	cw.FormatCaller = func(caller interface{}) string {
		if caller == nil || caller == "" {
			return ""
		}

		// 将 caller 字符串转换为 string 类型
		callerStr, ok := caller.(string)
		if !ok {
			return ""
		}

		// 转换为 "logger_test.go:76"
		return fmt.Sprintf("%s:", filepath.Base(callerStr))
	}
	return cw
}

// MultiLevelWriter may be used to send the log message to multiple outputs.
func MultiLevelWriter(writers ...io.Writer) zerolog.LevelWriter {
	return zerolog.MultiLevelWriter(writers...)
}

// 在 init 函数中设置全局的 Caller 格式化函数
func init() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		// 只返回文件名和行数，例如 "logger_test.go:76"
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
}

// New returns a new DefaultLogger instance
func New(options ...Opt) *DefaultLogger {
	var l = zerolog.New(os.Stdout).With().CallerWithSkipFrameCount(4).Logger()
	traceHookConfig := &TraceHookConfig{
		recordStackTraceInSpan: true,
		enableLevels:           AllLevel,
		errorSpanLevel:         zerolog.ErrorLevel,
	}
	// add request_id hook
	options = append(options, WithHookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		if e.GetCtx() == nil {
			return
		}
		logId, ok := e.GetCtx().Value(ReqIDKey).(string)
		if ok {
			e.Str(LogIDKey, logId)
		}
	}))
	options = append(options, WithHook(NewTraceHook(traceHookConfig)))
	return newLogger(l, options)
}

// From returns a new DefaultLogger instance using an existing logger
func From(log zerolog.Logger, options ...Opt) *DefaultLogger {
	return newLogger(log, options)
}

// GetLogger returns the default logger instance
func GetLogger() (DefaultLogger, error) {
	defaultlogger := GetDefaultLogger()

	if defaultlogger != nil {
		return *defaultlogger, nil
	}
	//if l, ok := DefaultLogger.(*DefaultLogger); ok {
	//	return *l, nil
	//}

	return DefaultLogger{}, errors.New("GetDefaultLogger is not a zerolog logger")
}

// SetLevel setting logging level for logger
func (l *DefaultLogger) SetLevel(level Level) {
	lvl := matchHlogLevel(level)
	l.level = lvl
	l.log = l.log.Level(lvl)
}

// SetOutput setting output for logger
func (l *DefaultLogger) SetOutput(writer io.Writer) {
	l.out = writer
	l.log = l.log.Output(writer)
}

// WithContext returns context with logger attached
func (l *DefaultLogger) WithContext(ctx context.Context) context.Context {
	return l.log.WithContext(ctx)
}

// WithField appends a field to the logger
func (l *DefaultLogger) WithField(key string, value interface{}) DefaultLogger {
	l.log = l.log.With().Interface(key, value).Logger()
	return *l
}

// Unwrap returns the underlying zerolog logger
func (l *DefaultLogger) Unwrap() zerolog.Logger {
	return l.log
}

// Log log using zerolog logger with specified level
func (l *DefaultLogger) Log(level Level, kvs ...interface{}) {
	switch level {
	case LevelTrace, LevelDebug:
		l.log.Debug().Msg(fmt.Sprint(kvs...))
	case LevelInfo:
		l.log.Info().Msg(fmt.Sprint(kvs...))
	case LevelNotice, LevelWarn:
		l.log.Warn().Msg(fmt.Sprint(kvs...))
	case LevelError:
		l.log.Error().Msg(fmt.Sprint(kvs...))
	case LevelFatal:
		l.log.Fatal().Msg(fmt.Sprint(kvs...))
	default:
		l.log.Warn().Msg(fmt.Sprint(kvs...))
	}
}

// Logf log using zerolog logger with specified level and formatting
func (l *DefaultLogger) Logf(level Level, format string, kvs ...interface{}) {
	switch level {
	case LevelTrace, LevelDebug:
		l.log.Debug().Msg(fmt.Sprintf(format, kvs...))
	case LevelInfo:
		l.log.Info().Msg(fmt.Sprintf(format, kvs...))
	case LevelNotice, LevelWarn:
		l.log.Warn().Msg(fmt.Sprintf(format, kvs...))
	case LevelError:
		l.log.Error().Msg(fmt.Sprintf(format, kvs...))
	case LevelFatal:
		l.log.Fatal().Msg(fmt.Sprintf(format, kvs...))
	default:
		l.log.Warn().Msg(fmt.Sprintf(format, kvs...))
	}
}

// CtxLogf log with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *DefaultLogger) CtxLogf(level Level, ctx context.Context, format string, kvs ...interface{}) {
	//logId, _ := ctx.Value(ReqIDKey).(string)

	unwrap := l.Unwrap()
	// todo add hook
	switch level {
	case LevelTrace, LevelDebug:
		unwrap.Debug().Ctx(ctx).Msg(fmt.Sprintf(format, kvs...))
	case LevelInfo:
		unwrap.Info().Ctx(ctx).Msg(fmt.Sprintf(format, kvs...))
	case LevelNotice, LevelWarn:
		unwrap.Warn().Ctx(ctx).Msg(fmt.Sprintf(format, kvs...))
	case LevelError:
		unwrap.Error().Ctx(ctx).Msg(fmt.Sprintf(format, kvs...))
	case LevelFatal:
		unwrap.Fatal().Ctx(ctx).Msg(fmt.Sprintf(format, kvs...))
	default:
		unwrap.Warn().Ctx(ctx).Msg(fmt.Sprintf(format, kvs...))
	}
}

// Trace logs a message at trace level.
func (l *DefaultLogger) Trace(v ...interface{}) {
	l.Log(LevelTrace, v...)
}

// Debug logs a message at debug level.
func (l *DefaultLogger) Debug(v ...interface{}) {
	l.Log(LevelDebug, v...)
}

// Info logs a message at info level.
func (l *DefaultLogger) Info(v ...interface{}) {
	l.Log(LevelInfo, v...)
}

// Notice logs a message at notice level.
func (l *DefaultLogger) Notice(v ...interface{}) {
	l.Log(LevelNotice, v...)
}

// Warn logs a message at warn level.
func (l *DefaultLogger) Warn(v ...interface{}) {
	l.Log(LevelWarn, v...)
}

// Error logs a message at error level.
func (l *DefaultLogger) Error(v ...interface{}) {
	l.Log(LevelError, v...)
}

// Fatal logs a message at fatal level.
func (l *DefaultLogger) Fatal(v ...interface{}) {
	l.Log(LevelFatal, v...)
}

// Tracef logs a formatted message at trace level.
func (l *DefaultLogger) Tracef(format string, v ...interface{}) {
	l.Logf(LevelTrace, format, v...)
}

// Debugf logs a formatted message at debug level.
func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	l.Logf(LevelDebug, format, v...)
}

// Infof logs a formatted message at info level.
func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	l.Logf(LevelInfo, format, v...)
}

// Noticef logs a formatted message at notice level.
func (l *DefaultLogger) Noticef(format string, v ...interface{}) {
	l.Logf(LevelWarn, format, v...)
}

// Warnf logs a formatted message at warn level.
func (l *DefaultLogger) Warnf(format string, v ...interface{}) {
	l.Logf(LevelWarn, format, v...)
}

// Errorf logs a formatted message at error level.
func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	l.Logf(LevelError, format, v...)
}

// Fatalf logs a formatted message at fatal level.
func (l *DefaultLogger) Fatalf(format string, v ...interface{}) {
	l.Logf(LevelError, format, v...)
}

// CtxTracef logs a message at trace level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *DefaultLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(LevelTrace, ctx, format, v...)
}

// CtxDebugf logs a message at debug level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *DefaultLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(LevelDebug, ctx, format, v...)
}

// CtxInfof logs a message at info level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *DefaultLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(LevelInfo, ctx, format, v...)
}

// CtxNoticef logs a message at notice level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *DefaultLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(LevelNotice, ctx, format, v...)
}

// CtxWarnf logs a message at warn level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *DefaultLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(LevelWarn, ctx, format, v...)
}

// CtxErrorf logs a message at error level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *DefaultLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(LevelError, ctx, format, v...)
}

// CtxFatalf logs a message at fatal level with logger associated with context.
// If no logger is associated, DefaultContextLogger is used, unless DefaultContextLogger is nil, in which case a disabled logger is used.
func (l *DefaultLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.CtxLogf(LevelFatal, ctx, format, v...)
}

func newLogger(log zerolog.Logger, options []Opt) *DefaultLogger {
	opts := newOptions(log, options)

	return &DefaultLogger{
		log:     opts.context.Logger(),
		out:     nil,
		level:   opts.level,
		options: options,
	}
}

var loggerMutex sync.Mutex

func (l *DefaultLogger) SetLogger(v interface{}) {
	// 检查输入是否为nil
	if v == nil {
		return
	}

	// 类型断言检查
	if l, ok := v.(*DefaultLogger); ok {
		// 添加并发保护
		loggerMutex.Lock()
		defer loggerMutex.Unlock()

		logger = l
		return
	}

	// 可选：添加类型不匹配的日志记录
	// log.Printf("SetLogger: expected *DefaultLogger, got %T", v)
}
