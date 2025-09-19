# Repository Guidelines

## Project Structure & Module Organization
The Go domain layer lives in `core/`, split into adapters, contracts (including the CloudEvents envelope helpers), runners, and unit-of-work interfaces; place new process logic here so cross-transport behavior stays cohesive. Transport-specific code sits under `transports/` (per protocol subfolder—HTTP baseline plus optional NATS behind a build tag), while reusable unit-of-work runners for other runtimes live in `uows/`. Storage adapters now include both in-memory and S3/MinIO variants (build tag `s3`). Reference SDK bindings in `sdk/` before publishing new APIs to keep language parity. Worked examples and quickstart binaries belong in `examples/`; `make build` emits artifacts into `bin/`. Long-form explanations and RFCs go in `docs/`—follow the structure used by `docs/uow-guidelines.md` when adding references.

## Build, Test, and Development Commands
Use `make build` for the inline example binary and `make clean` to reset the workspace. Run `make test` (or `go test ./...`) locally before every push; include `-run` filters when focusing on a single package and set `GOCACHE=$(pwd)/.gocache GOTOOLCHAIN=local` when toolchain downloads are blocked. Mirror Python coverage with `PYTHONPATH=sdk/python python3 -m unittest discover -s sdk/python/tests -p 'test_*.py'`. For quick iterations you can execute packages directly with `go run ./path/to/cmd`.

## Coding Style & Naming Conventions
Format Go code with `gofmt` or `goimports` (tabs for indentation, exported identifiers in PascalCase). Keep package names short and lower_snake case (`core/simplecontent`). When touching Python support code under `sdk/python` or `uows/python`, apply `black` and `isort` before committing. Prefer descriptive filenames such as `queue_worker.go` that mirror the transport or unit handled.

## Testing Guidelines
Place Go tests alongside implementation files as `*_test.go` using the standard library `testing` package. Name test functions `Test<Feature>` and table-driven cases `t.Run` with scenario labels. Exercise new transports or storage adapters through the example flows and assert failures as well as happy paths. Guard anything that depends on live infrastructure (e.g., `nats-server`, real S3/MinIO) behind build tags such as `//go:build nats && integration` or `//go:build s3 && integration` so `make test` stays hermetic.

## Commit & Pull Request Guidelines
Follow the existing log—short, imperative messages (`Add makefile`, `Implement runner`). Group related edits per commit and note breaking changes in the body. Pull requests should describe intent, link tracking issues, and list manual or automated verification (paste the `make test` output summary). Include configuration updates, screenshots, or CLI transcripts when they influence reviewers.

## Security & Configuration Notes
Review `docs/security.md` before introducing new transports or external dependencies. Document secrets and environment variables in the PR description, prefer local `.env` files that stay out of version control, and call out any new CloudEvents types, storage credentials (S3/MinIO), or external brokers (Kafka, NATS, SQS) so downstream consumers stay aligned.
