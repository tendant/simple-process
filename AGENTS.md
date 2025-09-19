# Repository Guidelines

## Project Structure & Module Organization
The Go domain layer lives in `core/`, split into adapters, contracts, runners, and unit-of-work helpers; place new process logic here so cross-transport behavior stays cohesive. Transport-specific code sits under `transports/` (per protocol subfolder), while reusable unit-of-work runners for other runtimes live in `uows/`. Reference SDK bindings in `sdk/` before publishing new APIs to keep language parity. Worked examples and quickstart binaries belong in `examples/`; `make build` emits artifacts into `bin/`. Long-form explanations and RFCs go in `docs/`—follow the structure used by `docs/uow-guidelines.md` when adding references.

## Build, Test, and Development Commands
Use `make build` for the inline example binary and `make clean` to reset the workspace. Run `make test` (or `go test ./...`) locally before every push; include `-run` filters when focusing on a single package. For quick iterations you can execute packages directly with `go run ./path/to/cmd`.

## Coding Style & Naming Conventions
Format Go code with `gofmt` or `goimports` (tabs for indentation, exported identifiers in PascalCase). Keep package names short and lower_snake case (`core/simplecontent`). When touching Python support code under `sdk/python` or `uows/python`, apply `black` and `isort` before committing. Prefer descriptive filenames such as `queue_worker.go` that mirror the transport or unit handled.

## Testing Guidelines
Place Go tests alongside implementation files as `*_test.go` using the standard library `testing` package. Name test functions `Test<Feature>` and table-driven cases `t.Run` with scenario labels. Exercise new transports through the example flows and assert failures as well as happy paths. If your change needs integration fixtures (e.g., Kafka), guard them with build tags so `make test` remains hermetic.

## Commit & Pull Request Guidelines
Follow the existing log—short, imperative messages (`Add makefile`, `Implement runner`). Group related edits per commit and note breaking changes in the body. Pull requests should describe intent, link tracking issues, and list manual or automated verification (paste the `make test` output summary). Include configuration updates, screenshots, or CLI transcripts when they influence reviewers.

## Security & Configuration Notes
Review `docs/security.md` before introducing new transports or external dependencies. Document secrets and environment variables in the PR description, and prefer local `.env` files that stay out of version control.
