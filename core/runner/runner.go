package runner

import (
	"context"

	"github.com/tendant/simple-process/core/adapters"
	"github.com/tendant/simple-process/core/contracts"
	"github.com/tendant/simple-process/core/uow"
)

// Runner is the interface for a UoW runner.
type Runner interface {
	Run(ctx context.Context, uow uow.UoW, job contracts.Job) (*contracts.Result, error)
}

// SyncRunner executes a UoW synchronously in the same process.
// It's suitable for fast, lightweight UoWs.
type SyncRunner struct{}

// NewSyncRunner creates a new SyncRunner.
func NewSyncRunner() *SyncRunner {
	return &SyncRunner{}
}

// Run executes the UoW's Process method directly.
func (r *SyncRunner) Run(ctx context.Context, uow uow.UoW, job contracts.Job) (*contracts.Result, error) {
	return uow.Process(ctx, job)
}

// AsyncRunner sends a UoW job to a message bus for asynchronous processing.
// It's suitable for long-running or resource-intensive UoWs.
type AsyncRunner struct {
	Bus adapters.Bus
}

// NewAsyncRunner creates a new AsyncRunner.
func NewAsyncRunner(bus adapters.Bus) *AsyncRunner {
	return &AsyncRunner{Bus: bus}
}

// Run publishes the job to the configured message bus.
// It does not wait for the UoW to complete and returns nil result and error.
func (r *AsyncRunner) Run(ctx context.Context, uow uow.UoW, job contracts.Job) (*contracts.Result, error) {
	err := r.Bus.Publish(ctx, job)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

