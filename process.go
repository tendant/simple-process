// Package process provides a simple framework for building file processing pipelines.
// It includes contracts for jobs and results, unit of work interfaces,
// and runners for both synchronous and asynchronous execution.
package process

import (
	"github.com/tendant/simple-process/pkg/adapters"
	"github.com/tendant/simple-process/pkg/contracts"
	"github.com/tendant/simple-process/pkg/runner"
	"github.com/tendant/simple-process/pkg/uow"
)

// Re-export core types for easier consumption

// Contracts
type (
	Job      = contracts.Job
	File     = contracts.File
	Blob     = contracts.Blob
	Return   = contracts.Return
	Result   = contracts.Result
	Artifact = contracts.Artifact
)

// UoW interface
type UoW = uow.UoW

// Runner interfaces and implementations
type (
	Runner      = runner.Runner
	SyncRunner  = runner.SyncRunner
	AsyncRunner = runner.AsyncRunner
)

// Adapter interfaces
type (
	Bus      = adapters.Bus
	Storage  = adapters.Storage
	Metadata = adapters.Metadata
	Logger   = adapters.Logger
	Tracer   = adapters.Tracer
	Span     = adapters.Span
)

// Constructor functions
var (
	NewSyncRunner  = runner.NewSyncRunner
	NewAsyncRunner = runner.NewAsyncRunner
)