package bootstrap

import (
	"fmt"
	"strings"

	appcfg "gpu-scheduler-platform/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg appcfg.LoggingConfig) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(strings.ToLower(strings.TrimSpace(cfg.Level)))); err != nil {
		return nil, fmt.Errorf("parse log level %q: %w", cfg.Level, err)
	}

	zcfg := zap.Config{
		Level:             level,
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          normalizeLogEncoding(cfg.Format),
		EncoderConfig:     newEncoderConfig(normalizeLogEncoding(cfg.Format)),
		OutputPaths:       cloneOrDefault(cfg.OutputPaths, []string{"stdout"}),
		ErrorOutputPaths:  cloneOrDefault(cfg.ErrorOutputPaths, []string{"stderr"}),
	}

	lg, err := zcfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}
	return lg, nil
}

func newEncoderConfig(format string) zapcore.EncoderConfig {
	enc := zap.NewProductionEncoderConfig()
	enc.TimeKey = "ts"
	enc.LevelKey = "level"
	enc.NameKey = "logger"
	enc.CallerKey = "caller"
	enc.MessageKey = "msg"
	enc.StacktraceKey = "stacktrace"

	enc.EncodeTime = zapcore.ISO8601TimeEncoder
	enc.EncodeDuration = zapcore.StringDurationEncoder

	if format == "console" {
		enc.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		enc.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return enc
}

func normalizeLogEncoding(format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "console":
		return "console"
	case "json":
		return "json"
	default:
		return "console"
	}
}

func cloneOrDefault(in []string, def []string) []string {
	if len(in) == 0 {
		out := make([]string, len(def))
		copy(out, def)
		return out
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}
