package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/tysenmeroniuk/jobqueue/internal/queue"
)

// Reaper periodically returns expired-lease jobs back to the queue.
// Idempotent — safe to run from every worker concurrently.
type Reaper struct {
	q        *queue.Queue
	interval time.Duration
}

func NewReaper(q *queue.Queue, interval time.Duration) *Reaper {
	return &Reaper{q: q, interval: interval}
}

// Run ticks on the configured interval until ctx is canceled.
func (r *Reaper) Run(ctx context.Context) error {
	t := time.NewTicker(r.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			n, err := r.q.Reap(ctx, 5*time.Second)
			if err != nil {
				slog.Error("reap failed", "err", err)
				continue
			}
			if n > 0 {
				slog.Info("reaped expired leases", "count", n)
			}
		}
	}
}
