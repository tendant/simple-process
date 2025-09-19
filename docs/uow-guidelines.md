# UoW Guidelines

This document provides guidelines for creating new Units of Work (UoWs).

## Principles

* **Stateless & Idempotent:** UoWs should be stateless and idempotent. This means that given the same input job, the UoW should always produce the same result.
* **Focused:** Each UoW should have a single, well-defined responsibility.
* **Zero-Copy I/O:** UoWs should use presigned URLs to access file content directly from blob storage, avoiding unnecessary data copying.

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

