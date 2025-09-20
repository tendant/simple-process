package uow

import (
	"context"

	"github.com/tendant/simple-process/pkg/contracts"
)

// UoW defines the interface for a unit of work.
// Each UoW is responsible for a specific file processing step.
type UoW interface {
	// Process executes the unit of work.
	// It takes a job as input and returns a result or an error.
	Process(ctx context.Context, job contracts.Job) (*contracts.Result, error)
}

