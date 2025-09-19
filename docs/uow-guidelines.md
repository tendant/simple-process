# UoW Guidelines

This document provides guidelines for creating new Units of Work (UoWs).

## Principles

* **Stateless & Idempotent:** UoWs should be stateless and idempotent. This means that given the same input job, the UoW should always produce the same result.
* **Focused:** Each UoW should have a single, well-defined responsibility.
* **Zero-Copy I/O:** UoWs should use presigned URLs to access file content directly from blob storage, avoiding unnecessary data copying.
* **Observable:** Return attribute patches and artifacts that downstream services can reason about. If publishing via CloudEvents, include enough context (e.g., `job_id`, `uow`) for tracing.

## Go UoWs

Go UoWs should implement the `uow.UoW` interface:

```go
package uow

import (
	"context"

	"github.com/tendant/simple-process/core/contracts"
)

// UoW defines the interface for a unit of work.
type UoW interface {
	// Process executes the unit of work.
	Process(ctx context.Context, job contracts.Job) (*contracts.Result, error)
}
```

When emitting artifacts, prefer deterministic locations (e.g., `artifacts/<fileID>.sha256`) so consumers can look them up by convention.

## Python UoWs

Python UoWs should use the `@uow` decorator:

```python
from simple_process_sdk.uow import uow

@uow("my_uow")
def my_uow(job):
    # ...
    return {
        # ...
}
```

Ensure Python workers remain streaming-friendly: download payloads via presigned URLs and upload artifacts through the SDK utilities once implemented.

## Testing

- Write unit tests for every UoW covering success and failure paths (`*_test.go` in Go, `tests/test_*.py` in Python).
- Integration scenarios that depend on external brokers (e.g., NATS) must be guarded behind build tags or pytest markers so the default test suite stays hermetic.
- Use the CloudEvents helpers (`core/contracts/cloudevent.go`) when fabricating transport messages in tests.
