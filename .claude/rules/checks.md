---
description: "Exact pre-commit commands for the go-vendored CI archetype, kept in sync with ci.yaml"
---

# Pre-commit Checks

This repo's `ci.yaml` is generated from `dx/ci-templates/go-vendored.yaml` via `dx ci sync` (check
for drift with `dx ci drift`). Run these before committing so CI passes on the first try:

```bash
test -z "$(gofmt -s -l ./*.go logqlmodel/ | tee /dev/stderr)"   # first-party files only; vendored packages keep upstream style
go mod tidy -diff                                               # fails on go.mod/go.sum drift; run `go mod tidy` to fix
go vet ./...
go build ./...
go test -race -count=1 ./...
```

There is no `golangci-lint` job for this repo — vendored upstream code doesn't conform to it. Every
PR and every push to `main` gates `logqlmodel/` coverage at >=80%; vendored code is tracked only.

The daily audit re-runs `scripts/sync-upstream.sh` and requires a byte-clean `git diff`, so never
let a formatter rewrite a vendored file — the comparison has no normalisation, and a single
reformatted line fails the `Upstream Sync Drift` job until it is reverted.
