package logger

import "context"

type noopLogger struct{}

func (noopLogger) Debug(context.Context, string, ...any) {}
func (noopLogger) Info(context.Context, string, ...any)  {}
func (noopLogger) Warn(context.Context, string, ...any)  {}
func (noopLogger) Error(context.Context, string, ...any) {}
