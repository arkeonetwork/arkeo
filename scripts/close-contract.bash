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

BIN="arkeod"
BIN_TX="arkeo"
USER="$1"
PROVIDER="$2"
CHAIN="$3"

PUBKEY_RAW=$($BIN keys show "$PROVIDER" -p | jq -r .key)
PUBKEY=$($BIN debug pubkey-raw "$PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

CLIENT_PUBKEY_RAW=$($BIN keys show "$USER" -p | jq -r .key)
CLIENT_PUBKEY=$($BIN debug pubkey-raw "$CLIENT_PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

$BIN tx $BIN_TX close-contract -y --from "$USER" --gas auto -- "$PUBKEY" "$CHAIN" "$CLIENT_PUBKEY"
