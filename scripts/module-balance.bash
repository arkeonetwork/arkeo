#!/bin/bash

if [ -z "$1" ]; then
	echo "No module supplied"
	exit 1
fi

PWD=$(dirname -- "${BASH_SOURCE[0]}")
BIN="mercuryd"
TOKEN="token"
MODULE="$1"

ADDRESS=$($BIN query auth module-accounts -o json | jq -r ".accounts[] | select(.name==\"$MODULE\") | .base_account.address")

$BIN query bank balances --denom $TOKEN -o json -- "$ADDRESS"
