# sample-app-go

[![CI](https://github.com/mruthyunjaya-lakkappanavar/sample-app-go/actions/workflows/ci.yml/badge.svg)](https://github.com/mruthyunjaya-lakkappanavar/sample-app-go/actions/workflows/ci.yml)
[![Release](https://github.com/mruthyunjaya-lakkappanavar/sample-app-go/actions/workflows/release.yml/badge.svg)](https://github.com/mruthyunjaya-lakkappanavar/sample-app-go/actions/workflows/release.yml)

> Sample Go application using [GitHub shared reusable workflows](https://github.com/mruthyunjaya-lakkappanavar/github-shared-workflows).

## Overview

This is a minimal Go HTTP server (using gorilla/mux) that demonstrates how a consumer repository can leverage centralized CI/CD pipelines with just ~15 lines of workflow YAML.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check — returns `{"status": "ok", "version": "x.y.z"}` |
| GET | `/api/greet?name=X` | Greeting — returns `{"message": "Hello, X!"}` |

## Local Development

```bash
# Install dependencies
go mod tidy

# Run the app
go run main.go

# Run tests
go test -v ./...

# Run linter (requires golangci-lint)
golangci-lint run
```

## CI/CD

This repo uses **reusable workflows** from `github-shared-workflows`:

- **CI** (`ci.yml`): Lint (golangci-lint) → Test (go test) → Security scan (Trivy)
- **Release** (`release.yml`): Semantic versioning → Changelog → GitHub Release → Slack notify

Both workflows are ~15 lines each — all logic lives in the central shared repo.

## License

MIT
