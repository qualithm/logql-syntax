.PHONY: build test test-race lint fmt sync

build:
	go build ./...

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	gofmt -s -w ./*.go logqlmodel/
	command -v goimports >/dev/null && goimports -w -local github.com/qualithm/logql-syntax ./*.go logqlmodel/ || true

lint:
	go vet ./...
	command -v golangci-lint >/dev/null && golangci-lint run --timeout=2m || echo "golangci-lint not installed; skipping"

# Re-sync vendored Loki source. See scripts/sync-upstream.sh for the version pin.
sync:
	./scripts/sync-upstream.sh
