# Repository Guidelines

## Project Layout
Go sources live at `main.go` for the CLI entry point and under `internal/cli` for command parsing, option handling, and reusable helpers. Unit tests (`cli_test.go`, `export_test.go`) and golden inputs in `internal/cli/testdata` cover the CLI package. Sample fixtures for documentation and demos sit in `demo/`. The compiled executable created locally defaults to `bin/purl`; a prebuilt binary used for demos lives at `purl/`.

## Build, Test & Tooling
The project targets Go 1.25; ensure your local toolchain matches `go.mod`. Common workflows:
- `make all` builds `bin/purl` with the embedded git revision.
- `make test` (or `go test ./...`) runs the suite with verbose output and coverage.
- `make vet` invokes `go vet` for static analysis, and `make staticcheck` runs `staticcheck -checks="all,-ST1000"`.
- `make cover` produces `coverage.out` and opens the HTML report via `go tool cover -html=coverage.out`.

## Coding Style
Run `goimports` (or `goimports ./...`) before committing so imports stay organized; if `goimports` is unavailable, fall back to `gofmt`. Go tooling expects tab indentation and CamelCase identifiers for exported symbols. Favor table-driven tests and keep option names aligned with CLI flags (e.g., `FilterFlag`). Avoid introducing external dependencies without discussing the impact on the single-binary distribution.

## Testing Expectations
Add or update `*_test.go` files alongside changed code. Table-driven tests in `internal/cli/cli_test.go` are the preferred pattern. Use fixtures in `internal/cli/testdata` or place new ones nearby with descriptive names (e.g., `filter-multi-line.txt`). Run `make test` before opening a PR; include `make cover` when changes touch parsing or filtering logic.

## Commits & Pull Requests
Commits should use concise, imperative subjects (e.g., `Add glob support to filter`) and group related changes. Reference issues when relevant using `Fixes #123`. For pull requests, provide a short summary, testing notes (`make test`, `make staticcheck`), and example command output or before/after snippets when UI or CLI behavior changes. Screenshots or asciinema links are encouraged for demo updates. Request at least one review before merging and ensure CI (if configured) is green.
