# File Processing Pipeline Library

A language-agnostic post-processing toolkit that turns each downstream operation—hashing, OCR, thumbnails, embeddings—into an independent Unit of Work (UoW). The Go core integrates with simplecontent, while SDKs and transports let you compose resilient pipelines that match your durability and orchestration requirements.

## Core Principles
- **Unit of Work first:** Every processor implements `uow.UoW` or an SDK decorator so steps stay isolated, idempotent, and retry-safe.
- **Bring your own orchestration:** Run inline, fan out over a queue, or plug into DBOS/user-defined DAGs; the library remains orchestration-agnostic.
- **Zero-copy I/O:** Workers stream data via presigned URLs instead of copying payloads through intermediaries.
- **Resume anywhere:** Persist state after every UoW and rely on idem keys to recover without double effects.
- **Pluggable adapters:** Swap storage, metadata, bus, logging, and tracing implementations without touching business logic.
- **Observability ready:** Runners stay lean so you can thread in tracing, logging, or metrics adapters without touching UoWs.

## Repository Layout
- `core/` – Go domain primitives (contracts, adapters, runners, UoW interface) that applications embed.
- `transports/` – Protocol-specific bindings (HTTP callback handler baseline plus an optional NATS bus under the `nats` build tag that emits CloudEvents envelopes).
- `examples/` – Runnable walkthroughs (currently inline hash example) showing how to wire a runner, adapter, and UoW.
- `uows/` – Reference UoWs in multiple languages (`uows/go`, `uows/python`) for reuse across services.
- `sdk/` – Language SDKs that expose decorators/helpers for registering UoWs with their runtimes.
- `docs/` – Design notes and contract references that complement this README.

## Getting Started
1. Install Go 1.18+ and (optionally) Python 3.10+ if you plan to run the Python SDK.
2. Build the inline sample: `make build` (outputs to `bin/inline-example`).
3. Execute tests: `make test` (or `go test ./...`) to cover unit and async integration scenarios. Set `GOCACHE=$(pwd)/.gocache` (and `GOTOOLCHAIN=local` when toolchain downloads are blocked) in sandboxed environments.
4. Run the inline example: `go run ./examples/inline` after pointing `storage.Put` at a reader for your input file.
5. Try the async workflow: `go run ./examples/async` to see `AsyncRunner` publishing to the in-memory bus while a worker updates metadata.
6. Validate the Python SDK: `PYTHONPATH=sdk/python python3 -m unittest discover -s sdk/python/tests -p 'test_*.py'`.
7. (Optional) Run the NATS demo once a local `nats-server` is running: `go run -tags nats ./examples/nats` (requires `go get github.com/nats-io/nats.go`). Jobs are wrapped in CloudEvents v1.0 envelopes, so any downstream consumer that speaks CloudEvents can participate.

## Working with Units of Work
- **Go:** Implement `core/uow.UoW` and return a `contracts.Result`. Refer to `uows/go/hash/hash.go` for a minimal example.
- **Python:** Decorate a function with `@uow("name")` from `sdk/python/uow`. Keep return payloads JSON-serializable and mirror the `Result` contract.
- Persist artifacts via the `adapters.Storage` interface and update metadata using `adapters.Metadata`.

## Choosing a Runner
- Use `core/runner.SyncRunner` for inline execution inside an API or CLI process.
- Use `core/runner.AsyncRunner` with an `adapters.Bus` implementation to fan jobs out to external workers.
- Compose runners with tracing/logging adapters so cross-cutting concerns stay outside UoW code.

## Embedding in Your Service
- Add the module: `go get github.com/tendant/simple-process@latest` for Go services, or install the Python SDK (`PYTHONPATH=sdk/python` during development) for worker code.
- Inject adapters that reflect your infrastructure (e.g., S3-backed storage, Dynamo metadata, Kafka/NATS bus) while keeping UoWs oblivious to deployment details.
- Register or import your UoWs (`uows/go/...`, `uows/python/...`) and execute them via `SyncRunner` (inline) or `AsyncRunner` (queue-based) depending on latency and durability needs.
- Persist the returned `contracts.Result` by patching metadata, recording artifacts, or chaining additional jobs; use transports/handlers to publish follow-up CloudEvents if required.
- Cover the workflow with tests: reuse the async example as an integration template and mirror the Python test command for multi-language validation.

## Extending the Library
- Implement additional transports (Kafka, NATS, SQS) under `transports/` by translating incoming jobs into `contracts.Job`.
- Provide concrete adapters in `core/adapters/*` to integrate with your blob store, metadata service, or observability stack; the in-memory bus and metadata adapters double as reference implementations.
- Add new reference UoWs under `uows/` and document them in `docs/` so other teams can reuse them.
- Keep Job/Result evolution backward compatible; document contract changes in `docs/contracts.md` and version payloads via the `Job.Version` field.

## NATS Queue Walkthrough (Optional)
- Start a local broker: `nats-server` (Homebrew: `brew install nats-server`).
- Fetch the NATS client once: `go get github.com/nats-io/nats.go@latest`.
- Publish and consume a job via NATS: `go run -tags nats ./examples/nats`. The example wires `AsyncRunner` into the NATS-backed bus and processes the message with a queue worker using the same in-memory storage used elsewhere in the repository while wrapping every message in a CloudEvents v1.0 envelope.

## CloudEvents Envelope
- Jobs published over transports are wrapped in a minimal CloudEvents v1.0 structure (`core/contracts/cloudevent.go`).
- The event `type` is `simpleprocess.job`, `id` mirrors `job_id`, and the payload lives in `data` with `datacontenttype` set to `application/json`.
- Consumers that already support CloudEvents can leverage headers for routing, retries, and schema enforcement without custom glue.

## Project Status & Next Steps
This repository is still a scaffold: storage adapters, transports, and SDK utilities are minimal. Before production use, flesh out real bus/metadata/logging implementations, complete end-to-end examples, and automate tests across supported languages.
