package worker

import (
	"context"

	"github.com/tysenmeroniuk/jobqueue/internal/handler"
	"github.com/tysenmeroniuk/jobqueue/internal/queue"
)

// Worker leases jobs from the queue, dispatches them to the registered
// handler for their type, and reports success/failure back to the queue.
type Worker struct {
	q        *queue.Queue
	registry *handler.Registry
	id       string
}

func New(q *queue.Queue, registry *handler.Registry, id string) *Worker {
	return &Worker{q: q, registry: registry, id: id}
}

// Run starts the lease loop. Blocks until ctx is canceled.
func (w *Worker) Run(ctx context.Context) error {
	// TODO: lease loop
	//   1. sleep/poll (or listen/notify) until a job is available
	//   2. q.Lease(ctx, w.id, queues, leaseDur)
	//   3. lookup handler by job.Type
	//   4. execute with per-job timeout ≤ lease duration
	//   5. on success: q.Complete; on failure: q.Retry (if attempts left) or q.Fail
	<-ctx.Done()
	return ctx.Err()
}
