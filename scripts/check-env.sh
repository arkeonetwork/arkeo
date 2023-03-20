#!/bin/bash

# skip checks if there is no tty (this is being called in a build script)
if [ -f /.dockerenv ] || [ ! -f /dev/tty ] || [ -t 0 ]; then
	exit 0
fi

set -euo pipefail

# ------------------------------ GNU Utilities ------------------------------

FAILED=0

check_gnu() {
	$1 --version 2>/dev/null | head -n 1 | grep -q "GNU" && return
	echo "GNU $1 is required." >/dev/tty
	FAILED=1
}

check_gnu grep
check_gnu awk
check_gnu find
check_gnu sed

# check make is version 4+
if [ -z "$FAILED" ]; then
	MAKE_VERSION=$(make --version 2>/dev/null | head -n 1 | awk -F '[ \\.]' '{print $3}')
	if [ "$MAKE_VERSION" -lt 4 ]; then
		echo "make version 4+ is required." >/dev/tty
		FAILED=1
	fi
fi

if [ "$FAILED" -gt 0 ]; then
	if [ "$(uname)" == "Darwin" ]; then
		echo >/dev/tty
		echo "Mac OS can try the following to update native utilities to the latest GNU version (homebrew):" >/dev/tty
		echo "1. Run: brew install coreutils binutils diffutils findutils gnu-tar gnu-sed gawk grep make" >/dev/tty
		echo "2. Follow the instructions from the brew output to update your PATH so the GNU utilities are default" >/dev/tty
	fi
	echo FAILED
fi

# ------------------------------ Golang ------------------------------

version() { echo "$@" | awk -F. '{ printf("%d%03d%03d%03d\n", $1,$2,$3,$4); }'; }

# Check Go version.
GO_VER=$(go version | grep -Eo 'go[0-9.]+' | sed -e s/go//)
MIN_VER="1.18.0"

# shellcheck disable=SC2046
if [ $(version "$GO_VER") -lt $(version "$MIN_VER") ]; then
	cat >/dev/tty <<EOF
Error: Detected Go version $GO_VER - this repository requires Go $MIN_VER as a minimum.
Please update Go and try again.
EOF
	echo FAILED
fi

# ------------------------------ Docker ------------------------------

if ! docker compose version >/dev/null; then
	echo "Docker and Compose plugin is required: https://docs.docker.com/compose/install/" >/dev/tty
	echo FAILED
fi
