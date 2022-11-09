#!/bin/bash

if [ -z "$1" ]; then
	echo "No user supplied"
	exit 1
fi

if [ -z "$2" ]; then
	echo "No provider supplied"
	exit 1
fi

if [ -z "$3" ]; then
	echo "No chain supplied"
	exit 1
fi

if [ -z "$4" ]; then
	echo "No contract type supplied"
	exit 1
fi

if [ -z "$5" ]; then
	echo "No duration supplied"
	exit 1
fi

BIN="mercuryd"
BIN_TX="mercury"
USER="$1"
PROVIDER="$2"
CHAIN="$3"
CTYPE="$4"
DURATION="$5"
RATE="10"

PUBKEY_RAW=$($BIN keys show "$PROVIDER" -p | jq -r .key)
PUBKEY=$($BIN debug pubkey-raw "$PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

CLIENT_PUBKEY_RAW=$($BIN keys show "$USER" -p | jq -r .key)
CLIENT_PUBKEY=$($BIN debug pubkey-raw "$CLIENT_PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

# nolint
DEPOSIT=$((RATE * DURATION))

$BIN tx $BIN_TX open-contract -y --from "$USER" --gas auto -- "$PUBKEY" "$CHAIN" "$CLIENT_PUBKEY" "$CTYPE" "$DEPOSIT" "$DURATION" $RATE
