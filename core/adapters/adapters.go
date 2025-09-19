package adapters

import (
	"context"
	"io"

	"github.com/tendant/simple-process/core/contracts"
)

// Storage provides an interface for blob storage operations.
type Storage interface {
	// Get returns a reader for the given blob location.
	Get(ctx context.Context, location string) (io.ReadCloser, error)
	// Put uploads a blob from a reader to the given location.
	Put(ctx context.Context, location string, reader io.Reader) error
	// PresignGet generates a presigned URL for getting a blob.
	PresignGet(ctx context.Context, location string) (string, error)
}

// Metadata provides an interface for interacting with file and artifact metadata.
type Metadata interface {
	// UpdateFileAttributes updates the attributes of a file.
	UpdateFileAttributes(ctx context.Context, fileID string, attributesPatch map[string]interface{}) error
	// CreateArtifact creates a new artifact record.
	CreateArtifact(ctx context.Context, fileID string, artifact contracts.Artifact) error
}

// Bus provides an interface for publishing jobs to a message bus.
type Bus interface {
	// Publish sends a job to the bus.
	Publish(ctx context.Context, job contracts.Job) error
}

// Logger provides a structured logging interface.
type Logger interface {
	// Info logs an informational message.
	Info(msg string, keysAndValues ...interface{})
	// Error logs an error message.
	Error(err error, msg string, keysAndValues ...interface{})
}

// Tracer provides an interface for creating spans for distributed tracing.
type Tracer interface {
	// StartSpan starts a new tracing span.
	StartSpan(ctx context.Context, name string) (context.Context, Span)
}

// Span is an interface for a single tracing span.
type Span interface {
	// End ends the span.
	End()
}

