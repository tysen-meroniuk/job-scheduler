package queue

// All SQL queries used by the queue live here as named constants so reviewers
// can read them in one place. The leasing query is the centerpiece.

const (
	// enqueueJob inserts a new job. If an idempotency_key is supplied and
	// collides with an existing job, DO NOTHING (caller resolves existing id).
	// $1=queue $2=type $3=payload $4=priority $5=run_at $6=max_attempts $7=idempotency_key
	enqueueJob = `
INSERT INTO jobs (queue, type, payload, priority, run_at, max_attempts, idempotency_key)
VALUES ($1, $2, $3, $4, COALESCE($5, NOW()), $6, $7)
ON CONFLICT (idempotency_key) DO NOTHING
RETURNING id
`

	// lookupByIdempotencyKey resolves the existing job id when enqueue collides.
	lookupByIdempotencyKey = `
SELECT id FROM jobs WHERE idempotency_key = $1
`

	// leaseOne is the heart of the system. FOR UPDATE SKIP LOCKED lets many
	// workers pull concurrently without contention. The partial index
	// idx_jobs_lease makes this fast regardless of table size.
	// $1=worker_id $2=lease_seconds $3=queues (text[])
	leaseOne = `
UPDATE jobs
SET state        = 'running',
    leased_by    = $1,
    leased_until = NOW() + ($2 * interval '1 second'),
    attempts     = attempts + 1,
    started_at   = COALESCE(started_at, NOW())
WHERE id = (
    SELECT id FROM jobs
    WHERE state = 'queued'
      AND run_at <= NOW()
      AND queue = ANY($3)
    ORDER BY priority DESC, run_at ASC
    FOR UPDATE SKIP LOCKED
    LIMIT 1
)
RETURNING id, queue, type, payload, state, priority, run_at,
          attempts, max_attempts, last_error, leased_by, leased_until,
          idempotency_key, created_at, started_at, completed_at
`

	// completeJob marks a successful execution. The leased_by check guards
	// against a worker completing a job whose lease has been reaped.
	// $1=job_id $2=worker_id
	completeJob = `
UPDATE jobs
SET state        = 'done',
    completed_at = NOW(),
    leased_by    = NULL,
    leased_until = NULL
WHERE id = $1 AND leased_by = $2
`

	// retryJob returns a failed-but-retryable job to queued with backoff.
	// $1=job_id $2=worker_id $3=backoff_seconds $4=error_message
	retryJob = `
UPDATE jobs
SET state        = 'queued',
    run_at       = NOW() + ($3 * interval '1 second'),
    last_error   = $4,
    leased_by    = NULL,
    leased_until = NULL
WHERE id = $1 AND leased_by = $2
`

	// failJob marks a job as permanently failed (attempts exhausted).
	// $1=job_id $2=worker_id $3=error_message
	failJob = `
UPDATE jobs
SET state        = 'failed',
    completed_at = NOW(),
    last_error   = $3,
    leased_by    = NULL,
    leased_until = NULL
WHERE id = $1 AND leased_by = $2
`

	// heartbeat extends a lease for a long-running job. The leased_by check
	// prevents a stale worker from extending a lease the reaper has already
	// reclaimed.
	// $1=job_id $2=worker_id $3=lease_seconds
	heartbeat = `
UPDATE jobs
SET leased_until = NOW() + ($3 * interval '1 second')
WHERE id = $1 AND leased_by = $2
`

	// reapExpired returns expired-lease jobs to queued state with a small
	// backoff. Does NOT increment attempts — attempts was already incremented
	// at lease time, so a crashed worker already cost the job one attempt
	// (poison-pill safety).
	// $1=backoff_seconds
	reapExpired = `
UPDATE jobs
SET state        = 'queued',
    run_at       = NOW() + ($1 * interval '1 second'),
    leased_by    = NULL,
    leased_until = NULL
WHERE state = 'running' AND leased_until < NOW()
RETURNING id
`
)
