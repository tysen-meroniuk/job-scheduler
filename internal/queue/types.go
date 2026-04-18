package queue

import (
	"encoding/json"
	"time"
)

type State string

const (
	StateQueued  State = "queued"
	StateRunning State = "running"
	StateDone    State = "done"
	StateFailed  State = "failed"
)

// Job is the canonical representation of a queued work item.
type Job struct {
	ID             int64           `json:"id"`
	Queue          string          `json:"queue"`
	Type           string          `json:"type"`
	Payload        json.RawMessage `json:"payload"`
	State          State           `json:"state"`
	Priority       int16           `json:"priority"`
	RunAt          time.Time       `json:"run_at"`
	Attempts       int             `json:"attempts"`
	MaxAttempts    int             `json:"max_attempts"`
	LastError      *string         `json:"last_error,omitempty"`
	LeasedBy       *string         `json:"leased_by,omitempty"`
	LeasedUntil    *time.Time      `json:"leased_until,omitempty"`
	IdempotencyKey *string         `json:"idempotency_key,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	StartedAt      *time.Time      `json:"started_at,omitempty"`
	CompletedAt    *time.Time      `json:"completed_at,omitempty"`
}
