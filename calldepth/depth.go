// Copyright (c) 2023 Breu Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package calldepth

import (
	"context"
	"log/slog"
	"runtime"
	"sync/atomic"
	"time"
)

type (
	// Adapter is the interface that wraps the slog.Logger with a call depth.
	Adapter interface {
		Debug(msg string, args ...any)
		DebugContext(ctx context.Context, msg string, args ...any)
		Enabled(ctx context.Context, level slog.Level) bool
		Error(msg string, args ...any)
		ErrorContext(ctx context.Context, msg string, args ...any)
		Handler() slog.Handler
		Info(msg string, args ...any)
		InfoContext(ctx context.Context, msg string, args ...any)
		Log(ctx context.Context, level slog.Level, msg string, args ...any)
		LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
		Warn(msg string, args ...any)
		WarnContext(ctx context.Context, msg string, args ...any)
		With(args ...any) Adapter
		WithGroup(name string) Adapter
	}

	// adapter is the implementation of Adapter.
	adapter struct {
		logger *slog.Logger // logger is the underlying logger.
		depth  int          // depth gives the call depth of the caller. DefaultCallDepth is 3. For 3rd pary adapters, this should be 4.
	}

	// Option provides a way to configure the adapter.
	Option func(*adapter)
)

const (
	// DefaultCallDepth helps skip Callers, the adapter.log function, and the adapter.log function's caller.
	DefaultCallDepth = 3
)

var (
	store atomic.Value
)

func (a *adapter) Enabled(ctx context.Context, level slog.Level) bool {
	return a.logger.Enabled(ctx, level)
}

func (a *adapter) Handler() slog.Handler {
	return a.logger.Handler()
}

func (a *adapter) With(args ...any) Adapter {
	return &adapter{
		logger: a.logger.With(args...),
		depth:  a.depth,
	}
}

func (a *adapter) WithGroup(name string) Adapter {
	return &adapter{
		logger: a.logger.WithGroup(name),
		depth:  a.depth,
	}
}

func (a *adapter) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	a.log(ctx, level, msg, args...)
}

func (a *adapter) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	a.logattrs(ctx, level, msg, attrs...)
}

func (a *adapter) Debug(msg string, args ...any) {
	a.log(context.Background(), slog.LevelDebug, msg, args...)
}

func (a *adapter) DebugContext(ctx context.Context, msg string, args ...any) {
	a.log(ctx, slog.LevelDebug, msg, args...)
}

func (a *adapter) Info(msg string, args ...any) {
	a.log(context.Background(), slog.LevelInfo, msg, args...)
}

func (a *adapter) InfoContext(ctx context.Context, msg string, args ...any) {
	a.log(ctx, slog.LevelInfo, msg, args...)
}

func (a *adapter) Warn(msg string, args ...any) {
	a.log(context.Background(), slog.LevelWarn, msg, args...)
}

func (a *adapter) WarnContext(ctx context.Context, msg string, args ...any) {
	a.log(ctx, slog.LevelWarn, msg, args...)
}

func (a *adapter) Error(msg string, args ...any) {
	a.log(context.Background(), slog.LevelError, msg, args...)
}

func (a *adapter) ErrorContext(ctx context.Context, msg string, args ...any) {
	a.log(ctx, slog.LevelError, msg, args...)
}

func (a *adapter) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if !a.logger.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(a.depth, pcs[:])

	record := slog.NewRecord(time.Now(), level, msg, pcs[0])

	record.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}

	_ = a.logger.Handler().Handle(ctx, record)
}

func (a *adapter) logattrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if !a.logger.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(a.depth, pcs[:])

	record := slog.NewRecord(time.Now(), level, msg, pcs[0])

	record.AddAttrs(attrs...)
	if ctx == nil {
		ctx = context.Background()
	}

	_ = a.logger.Handler().Handle(ctx, record)
}

func Default() Adapter {
	a := store.Load().(Adapter)
	if a == nil {
		a = New(WithSetDefault())
	}

	return a
}

func SetDefault(adapter Adapter) {
	store.Store(adapter)
}

// New returns a new Adapter.
func New(opts ...Option) Adapter {
	a := &adapter{
		logger: slog.Default(),
		depth:  DefaultCallDepth,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithLogger(logger *slog.Logger) Option {
	return func(a *adapter) {
		a.logger = logger
	}
}

func WithCallDepth(depth int) Option {
	return func(a *adapter) {
		a.depth = depth
	}
}

func WithSetDefault() Option {
	return func(a *adapter) {
		SetDefault(a)
	}
}
