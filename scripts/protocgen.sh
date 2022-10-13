#!/usr/bin/env bash

set -euo pipefail

# Delete any existing protobuf generated files.
find . -name "*.pb.go" -delete

go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos

# shellcheck disable=SC2038
find proto/ -path -prune -o -name '*.proto' -printf '%h\n' | sort | uniq |
	while read -r DIR; do
		find "$DIR" -maxdepth 1 -name '*.proto' |
			xargs protoc \
				-I "proto" \
				-I "third_party/proto" \
				--gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:.
	done

# Move proto files to the right places.
cp -r gitlab.com/thorchain/thornode/* ./
rm -rf gitlab.com
