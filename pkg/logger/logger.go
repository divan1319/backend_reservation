package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Environment string
	Level       string
	Rotation    RotationConfig
}

type RotationConfig struct {
	// Filename is the file to write logs to.
	Filename string
	// MaxSize is the maximum size in megabytes of the log file before it gets rotated.
	MaxSize int
	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int
	// MaxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename.
	MaxAge int
	// Compress determines if the rotated log files should be compressed using gzip.
	Compress bool
}

func InitLogger(cfg Config) {
	var handler slog.Handler

	logLevel := parseLevel(cfg.Level)

	switch strings.ToLower(cfg.Environment) {
	case "development", "dev":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: true, // Include source file and line number.
		})
	case "production", "prod":
		// For production, log to a file in JSON format with rotation.
		// This is ideal for log aggregation and analysis systems.
		logWriter := &lumberjack.Logger{
			Filename:   cfg.Rotation.Filename,
			MaxSize:    cfg.Rotation.MaxSize,
			MaxBackups: cfg.Rotation.MaxBackups,
			MaxAge:     cfg.Rotation.MaxAge,
			Compress:   cfg.Rotation.Compress,
		}

		// To ensure logs are written to both the file and console in production,
		// we use io.MultiWriter. Remove os.Stdout if you only want file logging.
		multiWriter := io.MultiWriter(os.Stdout, logWriter)

		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: true, // It's often useful to have source in production too.
		})

	default:
		// Default to a simple text handler if environment is not specified.
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Info("Logger successfully initialized", "environment", cfg.Environment)

}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error", "err":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type contextKey string

const loggerKey contextKey = "logger"

// CtxWithLogger creates a new context with a logger that includes the provided attributes.
// This is useful for adding request-specific context to logs.
func CtxWithLogger(ctx context.Context, attrs ...slog.Attr) context.Context {
	// Convert slog.Attr to []any for slog.With
	args := make([]any, 0, len(attrs)*2)
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value)
	}

	// Create a new logger with the provided attributes
	l := slog.Default().With(args...)

	// Store the logger in the context
	return context.WithValue(ctx, loggerKey, l)
}

// LoggerFromCtx retrieves the logger from the context.
// If no logger is found, it returns the default logger.
func LoggerFromCtx(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
