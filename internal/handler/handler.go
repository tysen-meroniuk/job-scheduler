package handler

import (
	"context"
	"fmt"
	"sync"
)

// Func is the signature every job handler implements. The worker invokes it
// with the job's payload; it returns an error if the job should retry/fail.
type Func func(ctx context.Context, payload []byte) error

// Registry maps job type strings to handlers.
type Registry struct {
	mu       sync.RWMutex
	handlers map[string]Func
}

func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string]Func)}
}

// Register associates a handler with a job type. Safe to call concurrently.
func (r *Registry) Register(jobType string, fn Func) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[jobType] = fn
}

// Lookup returns the handler for a job type, or an error if unregistered.
func (r *Registry) Lookup(jobType string) (Func, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn, ok := r.handlers[jobType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for job type %q", jobType)
	}
	return fn, nil
}
