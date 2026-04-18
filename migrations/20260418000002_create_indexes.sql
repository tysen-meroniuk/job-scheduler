-- +goose Up

-- The leasing index. Partial (only queued rows) so it stays tiny even as
-- done/failed rows accumulate. Column order matches the leasing query's
-- equality + ORDER BY, so Postgres can return rows without a sort step.
CREATE INDEX idx_jobs_lease
    ON jobs (queue, priority DESC, run_at ASC)
    WHERE state = 'queued';

-- Reaper index. Partial on running rows (typically a handful in flight),
-- keyed on leased_until so the expired-lease scan is a tiny range lookup.
CREATE INDEX idx_jobs_running_lease
    ON jobs (leased_until)
    WHERE state = 'running';

-- Idempotency. Partial unique index so NULL keys don't collide and don't
-- bloat the index. Enables INSERT ... ON CONFLICT (idempotency_key) DO NOTHING.
CREATE UNIQUE INDEX idx_jobs_idempotency
    ON jobs (idempotency_key)
    WHERE idempotency_key IS NOT NULL;

-- +goose Down
DROP INDEX idx_jobs_idempotency;
DROP INDEX idx_jobs_running_lease;
DROP INDEX idx_jobs_lease;
