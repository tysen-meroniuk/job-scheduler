package queue

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Queue is the entry point for all job queue operations.
type Queue struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Queue {
	return &Queue{pool: pool}
}

// EnqueueParams carries the inputs to Enqueue.
type EnqueueParams struct {
	Queue          string
	Type           string
	Payload        []byte
	Priority       int16
	RunAt          *time.Time // nil = run ASAP
	MaxAttempts    int        // 0 = use default (5)
	IdempotencyKey *string
}

// Enqueue inserts a new job. Returns the job id. If an idempotency key is
// supplied and matches an existing job, returns that job's id instead (no-op).
func (q *Queue) Enqueue(ctx context.Context, p EnqueueParams) (int64, error) {
	// TODO: implement using enqueueJob + lookupByIdempotencyKey on conflict
	return 0, nil
}

// Lease claims the highest-priority available job for the given queues.
// Returns nil, nil when nothing is available (not an error).
func (q *Queue) Lease(ctx context.Context, workerID string, queues []string, lease time.Duration) (*Job, error) {
	// TODO: implement using leaseOne
	return nil, nil
}

// Complete marks a job as successfully executed.
func (q *Queue) Complete(ctx context.Context, jobID int64, workerID string) error {
	// TODO: implement using completeJob
	return nil
}

// Retry returns a job to the queue with backoff (attempts remaining).
func (q *Queue) Retry(ctx context.Context, jobID int64, workerID string, errMsg string, backoff time.Duration) error {
	// TODO: implement using retryJob
	return nil
}

// Fail marks a job as permanently failed (no more attempts).
func (q *Queue) Fail(ctx context.Context, jobID int64, workerID string, errMsg string) error {
	// TODO: implement using failJob
	return nil
}

// Heartbeat extends the lease on an in-flight job.
func (q *Queue) Heartbeat(ctx context.Context, jobID int64, workerID string, lease time.Duration) error {
	// TODO: implement using heartbeat
	return nil
}

// Reap returns expired-lease jobs to the queue. Idempotent — safe to run
// from every worker on a timer. Returns the count of jobs reaped.
func (q *Queue) Reap(ctx context.Context, backoff time.Duration) (int, error) {
	// TODO: implement using reapExpired
	return 0, nil
}
