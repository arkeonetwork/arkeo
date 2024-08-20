#!/usr/bin/env bash

set -eo pipefail

cd proto

proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
	for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
		echo "checking $file"
		if grep go_package "$file" &>/dev/null; then
			echo "generating go files for $file"
			buf generate --template buf.gen.gogo.yaml "$file"
		fi
	done
done

cd ..

# after the proto files have been generated add them to the the repo
# in the proper location. Then, remove the ephemeral tree used for generation
echo "copying generated type files to repo"
sudo cp -r github.com/arkeonetwork/arkeo/*  .
rm -rf github.com

# we need to go mod manually, because the docker image is still on go1.18
# go mod tidy
