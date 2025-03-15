#/usr/bin/env bash

set -euxo pipefail

PYTHON_VERSION="3.10"
archive_path="${1:-mjai-app.zip}"

tempdir="$(mktemp -d)"
trap "rm -rf '$tempdir'" EXIT

cp test/mjai.app/test.py "$tempdir"

pushd "$tempdir"
uv venv --python $PYTHON_VERSION
uv pip install mjai
popd

while true; do
    logs_dir="./logs.$(date +%Y-%m-%d-%H-%M-%S)"
    uv run python "$tempdir/test.py" "$logs_dir" "$archive_path" || true
    grep -Fqr '"error"' "$logs_dir" || rm -rf "$logs_dir"
done
