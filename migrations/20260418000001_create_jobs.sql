-- +goose Up
CREATE TABLE jobs (
    id              BIGSERIAL      PRIMARY KEY,
    queue           TEXT           NOT NULL DEFAULT 'default',
    type            TEXT           NOT NULL,
    payload         JSONB          NOT NULL DEFAULT '{}'::jsonb,

    state           TEXT           NOT NULL DEFAULT 'queued',
    priority        SMALLINT       NOT NULL DEFAULT 0,
    run_at          TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    attempts        INT            NOT NULL DEFAULT 0,
    max_attempts    INT            NOT NULL DEFAULT 5,
    last_error      TEXT,

    leased_by       TEXT,
    leased_until    TIMESTAMPTZ,

    idempotency_key TEXT,

    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,

    CONSTRAINT jobs_state_check
        CHECK (state IN ('queued', 'running', 'done', 'failed'))
);

-- +goose Down
DROP TABLE jobs;
