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

BIN="mercuryd"
BIN_TX="mercury"
USER="$1"
CHAIN="$2"
BOND="$3"

PUBKEY_RAW=$($BIN keys show "$USER" -p | jq -r .key)
PUBKEY=$($BIN debug pubkey-raw "$PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

$BIN tx $BIN_TX bond-provider -y --from "$USER" --gas auto -- "$PUBKEY" "$CHAIN" "$BOND"
