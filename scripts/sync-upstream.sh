#!/usr/bin/env bash
# Re-sync vendored Loki LogQL source files into this fork.
#
# Usage:
#   scripts/sync-upstream.sh [LOKI_VERSION]
#
# Defaults to the version below. The script does not commit; review the diff
# manually, since `pkg/util` and `pkg/logqlmodel` are trimmed reimplementations
# that may need reconciliation against new upstream symbols.

set -euo pipefail

LOKI_VERSION="${1:-v3.7.2}"
MODULE_PATH="github.com/qualithm/logql-syntax"

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$repo_root"

# Resolve the Loki source in the local module cache.
GOPATH="$(go env GOPATH)"
loki="$GOPATH/pkg/mod/github.com/grafana/loki/v3@${LOKI_VERSION}"
if [[ ! -d "$loki" ]]; then
  echo "Loki ${LOKI_VERSION} is not in the module cache. Run:" >&2
  echo "  go mod download github.com/grafana/loki/v3@${LOKI_VERSION}" >&2
  exit 1
fi

# Copy vendored directories verbatim.
copy_dir() {
  local src="$1" dst="$2"
  rm -rf "$dst"
  mkdir -p "$dst"
  cp -R "$src"/. "$dst"/
}

copy_dir "$loki/pkg/logql/syntax"        "syntax"
copy_dir "$loki/pkg/logql/log"           "log"
copy_dir "$loki/pkg/logql/log/jsonexpr"  "log/jsonexpr"
copy_dir "$loki/pkg/logql/log/logfmt"    "log/logfmt"
copy_dir "$loki/pkg/logql/log/pattern"   "log/pattern"

# The vendored `log/` dir overwrites the subdirs above when re-copied. Restore.
copy_dir "$loki/pkg/logql/log/jsonexpr"  "log/jsonexpr"
copy_dir "$loki/pkg/logql/log/logfmt"    "log/logfmt"
copy_dir "$loki/pkg/logql/log/pattern"   "log/pattern"

# Re-copy the small trimmed util helpers (verbatim from upstream).
mkdir -p internal/util internal/util/encoding internal/constants
cp "$loki/pkg/util/regex.go"           internal/util/regex.go
cp "$loki/pkg/util/matchers.go"        internal/util/matchers.go
cp "$loki/pkg/util/encoding/encoding.go" internal/util/encoding/encoding.go
cp "$loki/pkg/util/constants/variants.go" internal/constants/variants.go

# Rewrite imports.
rewrite() {
  find "$@" -name '*.go' -print0 | xargs -0 sed -i.bak \
    -e "s|github.com/grafana/loki/v3/pkg/logql/syntax|${MODULE_PATH}/syntax|g" \
    -e "s|github.com/grafana/loki/v3/pkg/logql/log/jsonexpr|${MODULE_PATH}/log/jsonexpr|g" \
    -e "s|github.com/grafana/loki/v3/pkg/logql/log/logfmt|${MODULE_PATH}/log/logfmt|g" \
    -e "s|github.com/grafana/loki/v3/pkg/logql/log/pattern|${MODULE_PATH}/log/pattern|g" \
    -e "s|github.com/grafana/loki/v3/pkg/logql/log|${MODULE_PATH}/log|g" \
    -e "s|github.com/grafana/loki/v3/pkg/logqlmodel|${MODULE_PATH}/logqlmodel|g" \
    -e "s|github.com/grafana/loki/v3/pkg/util/encoding|${MODULE_PATH}/internal/util/encoding|g" \
    -e "s|github.com/grafana/loki/v3/pkg/util/constants|${MODULE_PATH}/internal/constants|g" \
    -e "s|github.com/grafana/loki/v3/pkg/util|${MODULE_PATH}/internal/util|g"
  find "$@" -name '*.go.bak' -delete
}
rewrite syntax log internal

# Refresh LICENSE from upstream.
cp "$loki/LICENSE" LICENSE

# Tidy and report.
go mod tidy
echo
echo "Sync complete against Loki ${LOKI_VERSION}."
echo "Review the diff, run \`go test ./...\`, and if logqlmodel/ needs new"
echo "symbols (see logqlmodel/logqlmodel.go) update it manually."
