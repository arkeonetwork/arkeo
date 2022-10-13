#!/usr/bin/env bash
set -euo pipefail

die() {
	echo "ERR: $*"
	exit 1
}

# Check that no .pb.go files were added.
if git ls-files '*.go' | grep -q '.pb.go$'; then
	die "Do not add generated protobuf .pb.go files"
fi

if [ -n "$(git ls-files '*.go' | grep -v -E '.pb.gw.go' | grep -v -e '^docs/' | xargs gofumpt -l 2>/dev/null)" ]; then
	git ls-files '*.go' | grep -v -E '.pb.gw.go' | grep -v -e '^docs/' | xargs gofumpt -d 2>/dev/null
	die "Go formatting errors"
fi

go mod verify
