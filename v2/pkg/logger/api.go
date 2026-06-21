package logger

import "context"

// Debug logs at debug level via the active process logger.
func Debug(ctx context.Context, msg string, args ...any) {
	active.Debug(ctx, msg, args...)
}

// Info logs at info level via the active process logger.
func Info(ctx context.Context, msg string, args ...any) {
	active.Info(ctx, msg, args...)
}

// Warn logs at warn level via the active process logger.
func Warn(ctx context.Context, msg string, args ...any) {
	active.Warn(ctx, msg, args...)
}

// Error logs at error level via the active process logger.
func Error(ctx context.Context, msg string, args ...any) {
	active.Error(ctx, msg, args...)
}
