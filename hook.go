package main

import (
	"errors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	logSeverityTextKey = attribute.Key("otel.log.severity.text")
	logMessageKey      = attribute.Key("otel.log.message")
)

var _ zerolog.Hook = (*TraceHook)(nil)

var AllLevel = []zerolog.Level{zerolog.TraceLevel, zerolog.DebugLevel, zerolog.InfoLevel,
	zerolog.WarnLevel, zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel}

type TraceHookConfig struct {
	recordStackTraceInSpan bool
	enableLevels           []zerolog.Level
	errorSpanLevel         zerolog.Level
}

type TraceHook struct {
	cfg *TraceHookConfig
}

func NewTraceHook(cfg *TraceHookConfig) *TraceHook {
	return &TraceHook{cfg: cfg}
}

func (h *TraceHook) Levels() []zerolog.Level {
	return h.cfg.enableLevels
}

func (h *TraceHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	if e.GetCtx() == nil {
		return
	}

	span := trace.SpanFromContext(e.GetCtx())
	if !span.IsRecording() {
		return
	}

	e.Str(traceIDKey, span.SpanContext().TraceID().String())
	e.Str(spanIDKey, span.SpanContext().SpanID().String())
	e.Str(traceFlagsKey, span.SpanContext().TraceFlags().String())

	// attach log to span event attributes
	attrs := []attribute.KeyValue{
		logMessageKey.String(message),
		logSeverityTextKey.String(OtelSeverityText(level)),
	}
	span.AddEvent(logEventKey, trace.WithAttributes(attrs...))

	// set span status
	if level >= h.cfg.errorSpanLevel {
		span.SetStatus(codes.Error, message)
		span.RecordError(errors.New(message), trace.WithStackTrace(h.cfg.recordStackTraceInSpan))
	}

	return
}

func OtelSeverityText(lv zerolog.Level) string {
	s := lv.String()
	//if s == "warning" {
	//	s = "warn"
	//}
	return strings.ToUpper(s)
}
