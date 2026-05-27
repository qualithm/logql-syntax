# LogQL Syntax

[![CI](https://github.com/qualithm/logql-syntax/actions/workflows/ci.yaml/badge.svg)](https://github.com/qualithm/logql-syntax/actions/workflows/ci.yaml)
[![codecov](https://codecov.io/gh/qualithm/logql-syntax/graph/badge.svg)](https://codecov.io/gh/qualithm/logql-syntax)
[![Go Reference](https://pkg.go.dev/badge/github.com/qualithm/logql-syntax.svg)](https://pkg.go.dev/github.com/qualithm/logql-syntax)
[![Go Report Card](https://goreportcard.com/badge/github.com/qualithm/logql-syntax)](https://goreportcard.com/report/github.com/qualithm/logql-syntax)

Standalone Go parser and AST for [Grafana Loki](https://github.com/grafana/loki)'s LogQL. Lifts the
upstream `syntax`, `log`, and `logqlmodel` packages out of `grafana/loki` with their runtime
dependencies (`dskit`, `etcd`, `jaeger`, the queryrange/push machinery) stripped away.

## Installation

```bash
go get github.com/qualithm/logql-syntax
```

## Usage

```go
import "github.com/qualithm/logql-syntax/syntax"

expr, err := syntax.ParseExpr(`sum by (job) (rate({app="api"} |= "error" [5m]))`)
if err != nil {
    return err
}
expr.Walk(func(e syntax.Expr) bool {
    // inspect the AST
    return true
})
```

## What's included

| Path             | Source                                                              |
| ---------------- | ------------------------------------------------------------------- |
| `syntax/`        | `github.com/grafana/loki/v3/pkg/logql/syntax`                       |
| `log/`           | `github.com/grafana/loki/v3/pkg/logql/log`                          |
| `log/jsonexpr/`  | `github.com/grafana/loki/v3/pkg/logql/log/jsonexpr`                 |
| `log/logfmt/`    | `github.com/grafana/loki/v3/pkg/logql/log/logfmt`                   |
| `log/pattern/`   | `github.com/grafana/loki/v3/pkg/logql/log/pattern`                  |
| `logqlmodel/`    | trimmed extract of `pkg/logqlmodel` (errors + label constants only) |
| `internal/util/` | three regex / matcher helpers from `pkg/util`                       |

The runtime `Result` and `Streams` types from `logqlmodel` are intentionally omitted because they
pull in `loki/pkg/push` and queryrange machinery.

## Upstream sync

Tracked against Loki [`v3.7.2`](https://github.com/grafana/loki/releases/tag/v3.7.2).

To resync against a newer Loki release:

1. Copy source files from `pkg/logql/{syntax,log}/...` into the matching directories here.
2. Rewrite `github.com/grafana/loki/v3/pkg/...` import paths to
   `github.com/qualithm/logql-syntax/...` (see the `sed` invocation in the project history).
3. Reconcile any new uses of `pkg/util`, `pkg/logqlmodel`, or `pkg/util/constants` — extend the
   trimmed packages here as needed.
4. `go test ./...` — the only known persistent failures are the two timestamp subtests in `log/`
   that hardcode local-timezone dates upstream.

## Development

### Prerequisites

- [Go](https://go.dev/dl/) 1.26+

### Setup

```bash
make install-tools
```

This installs local development tooling, including `golangci-lint`, `goimports`, `govulncheck`, and
`gosec`.

### Building & Testing

```bash
make build
make test
make lint
```

### Security Tooling

```bash
make audit   # govulncheck
make gosec   # standalone gosec scan
```

Daily CI security audit runs both tools in `.github/workflows/audit.yaml`.

Install tools manually (if you are not using `make install-tools`):

```bash
go install golang.org/x/vuln/cmd/govulncheck@v1.3.0
go install github.com/securego/gosec/v2/cmd/gosec@v2.26.1
```

### Resyncing from upstream Loki

```bash
make sync                       # defaults to LOKI_VERSION in the Makefile
LOKI_VERSION=v3.8.0 make sync   # pin a specific upstream release
```

## Minimum Supported Go Version

Go 1.26+.

## Licence

Apache-2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE) for upstream attribution.
