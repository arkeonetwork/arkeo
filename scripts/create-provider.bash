#!/bin/bash

if [ -z "$1" ]; then
	echo "No user supplied"
	exit 1
fi

if [ -z "$2" ]; then
	echo "No chain supplied"
	exit 1
fi

if [ -z "$3" ]; then
	echo "No bond supplied"
	exit 1
fi

PWD=$(dirname -- "${BASH_SOURCE[0]}")
BIN="arkeod"
BIN_TX="arkeo"
USER="$1"
CHAIN="$2"
BOND="$3"

PUBKEY_RAW=$($BIN keys show "$USER" -p | jq -r .key)
PUBKEY=$($BIN debug pubkey-raw "$PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

./"$PWD"/bond-provider.bash "$USER" "$CHAIN" "$BOND"

$BIN tx $BIN_TX mod-provider -y --from "$USER" --gas auto -- "$PUBKEY" "$CHAIN" "http://localhost:3636/metadata.json" 1 1 5 50 10 10
