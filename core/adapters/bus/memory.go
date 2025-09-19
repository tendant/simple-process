package bus

import (
	"context"

	"github.com/tendant/simple-process/core/contracts"
)

// MemoryBus is an in-memory implementation of adapters.Bus for examples and tests.
type MemoryBus struct {
	jobs chan contracts.Job
}

// NewMemoryBus creates a MemoryBus with the provided channel buffer size.
func NewMemoryBus(buffer int) *MemoryBus {
	if buffer <= 0 {
		buffer = 1
	}
	return &MemoryBus{jobs: make(chan contracts.Job, buffer)}
}

// Publish enqueues the job onto the in-memory buffer or returns ctx error on cancellation.
func (b *MemoryBus) Publish(ctx context.Context, job contracts.Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case b.jobs <- job:
		return nil
	}
}

// Subscribe exposes a read-only channel of published jobs for consumer goroutines.
func (b *MemoryBus) Subscribe() <-chan contracts.Job {
	return b.jobs
}

// Close closes the underlying channel, signalling consumers to stop.
func (b *MemoryBus) Close() {
	close(b.jobs)
}
