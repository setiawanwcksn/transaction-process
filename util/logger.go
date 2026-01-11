package util

import (
	"context"
	"io"
	"log/slog"
	"runtime"
	"time"
)

var Log *slog.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func Perf(ctx context.Context, label string, extra ...any) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	attrs := []any{
		"label", label,
		"ts", time.Now().Unix(),
		"goroutines", runtime.NumGoroutine(),
		"alloc_mb", m.Alloc / 1024 / 1024,
		"total_alloc_mb", m.TotalAlloc / 1024 / 1024,
		"sys_mb", m.Sys / 1024 / 1024,
		"num_gc", m.NumGC,
	}
	if len(extra) > 0 {
		attrs = append(attrs, extra...)
	}

	Log.InfoContext(ctx, "perf", attrs...)
}
