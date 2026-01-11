package worker

import (
	"context"
	"time"

	"flip/internal/model"
	"flip/util"
)

func (t *tracker) mark(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.seen[key] {
		return false
	}
	t.seen[key] = true
	return true
}

func (w *Repository) Start(ctx context.Context, workers int) {
	util.Log.InfoContext(ctx, "worker start", "threads", workers)

	for i := 0; i < workers; i++ {
		go w.loop(ctx, i)
	}
}

func (w *Repository) loop(ctx context.Context, wid int) {
	for {
		select {
		case <-w.stop:
			return
		case ev := <-w.bus:
			key := ev.ID

			if !w.tracker.mark(key) {
				continue
			}

			switch ev.Tx.Status {
			case model.StatusFailed:
				// w.handleFailed(ctx, ev, wid)
			case model.StatusPending:
				w.handlePending(ctx, ev, wid)
			default:
				// ignore success or unknown
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *Repository) handleFailed(ctx context.Context, ev model.Event, wid int) {
	defer func(start time.Time) {
		util.Perf(ctx, "worker_failed", "upload_id", ev.UploadID, "worker", wid, "elapsed_ms", time.Since(start).Milliseconds())
	}(time.Now())

	_ = w.handler(ctx, ev) // handler must call AddFailed/AddSuccess itself
}

func (w *Repository) handlePending(ctx context.Context, ev model.Event, wid int) {
	defer func(start time.Time) {
		util.Perf(ctx, "worker_pending", "upload_id", ev.UploadID, "worker", wid, "elapsed_ms", time.Since(start).Milliseconds())
	}(time.Now())

	w.sleep(2 * time.Second)

	for attempt := 1; attempt <= 3; attempt++ {
		res := w.handler(ctx, ev)
		if res.Success {
			return
		}
		if !res.Retryable {
			return
		}
		w.sleep(time.Duration(200*attempt) * time.Millisecond)
	}
}

func (w *Repository) Stop() {
	close(w.stop)
}
