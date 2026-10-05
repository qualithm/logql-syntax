# LogQL Syntax

[![CI](https://github.com/qualithm/logql-syntax/actions/workflows/ci.yaml/badge.svg)](https://github.com/qualithm/logql-syntax/actions/workflows/ci.yaml)
[![codecov](https://codecov.io/gh/qualithm/logql-syntax/graph/badge.svg)](https://codecov.io/gh/qualithm/logql-syntax)
[![Go Reference](https://pkg.go.dev/badge/github.com/qualithm/logql-syntax.svg)](https://pkg.go.dev/github.com/qualithm/logql-syntax)

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

| Path                  | Source                                                              |
| --------------------- | ------------------------------------------------------------------- |
| `syntax/`             | `github.com/grafana/loki/v3/pkg/logql/syntax`                       |
| `log/`                | `github.com/grafana/loki/v3/pkg/logql/log`                          |
| `log/jsonexpr/`       | `github.com/grafana/loki/v3/pkg/logql/log/jsonexpr`                 |
| `log/logfmt/`         | `github.com/grafana/loki/v3/pkg/logql/log/logfmt`                   |
| `log/pattern/`        | `github.com/grafana/loki/v3/pkg/logql/log/pattern`                  |
| `logqlmodel/`         | trimmed extract of `pkg/logqlmodel` (errors + label constants only) |
| `internal/util/`      | regex, matcher and encoding helpers from `pkg/util`                 |
| `internal/constants/` | `variants.go` from `pkg/util/constants`                             |

The runtime `Result` and `Streams` types from `logqlmodel` are intentionally omitted because they
pull in `loki/pkg/push` and queryrange machinery.

## Upstream sync

Tracked against Loki [`v3.7.2`](https://github.com/grafana/loki/releases/tag/v3.7.2).

To resync against a newer Loki release:

1. Download the release and run the sync script, which copies the upstream packages and rewrites
   their import paths. `make sync` uses the version pinned in `scripts/sync-upstream.sh`.

   ```bash
   go mod download github.com/grafana/loki/v3@<loki-version>
   ./scripts/sync-upstream.sh <loki-version>
   ```

2. Reconcile any new uses of `pkg/logqlmodel` — extend the trimmed `logqlmodel/` package here as
   needed.
3. `go test ./...` — the only known persistent failures are the two timestamp subtests in `log/`
   that hardcode local-timezone dates upstream.

## Development

### Prerequisites

- [Go](https://go.dev/dl/) 1.26+

### Setup

```bash
make install-tools
```

This installs `golangci-lint`, `goimports`, `govulncheck` and `gosec` into `$GOPATH/bin` (`~/go/bin`
by default). Put that directory on your `PATH`:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

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

`.github/workflows/audit.yaml` runs `govulncheck` and `gosec` daily.

## Minimum Supported Go Version

Go 1.26+.

## License

Apache-2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE) for upstream attribution.
