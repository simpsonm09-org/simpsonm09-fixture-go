# simpsonm09-fixture-go

A Go item CRUD service built on Gin. It is the Go fixture for the simpsonm09 fleet standard.

The original lives in `simpsonm09-org/simpsonm09-fixture-go`; work happens on the personal fork. See [`repo-standard`](https://github.com/simpsonm09-org/simpsonm09-repo-standard).

## What it does

The service exposes item CRUD over HTTP behind a layered architecture. A Gin handler speaks DTOs, a service works in the domain `Item`, and an in-memory store adapter implements a repository port. An item has a server-assigned id, a required name, and an optional description. The store is in memory, so restarting the process discards every item and restores three dev seeds.

## Quick start

```bash
just install
just verify
```

## Commands

| Command | Does |
| --- | --- |
| `just install` | Installs the pinned tools with `mise`. |
| `just lint` | Runs the linters. |
| `just test` | Runs the Go test suite. |
| `just coverage` | Runs the tests and writes `coverage/lcov.info`. |
| `just spec` | Regenerates `docs/openapi.json` from the handler annotations. |
| `just serve` | Serves the API on port 3000. |
| `just verify` | Lints and tests. |

## Documentation

Read [`docs/README.md`](docs/README.md) for the architecture, the items feature, and the OpenAPI contract.

## License

MIT. See [`LICENSE`](LICENSE).
