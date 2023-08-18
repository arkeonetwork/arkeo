#!/bin/bash

if [ -z "$1" ]; then
	echo "No user supplied"
	exit 1
fi

if [ -z "$2" ]; then
	echo "No service supplied"
	exit 1
fi

if [ -z "$3" ]; then
	echo "No bond supplied"
	exit 1
fi

BIN="arkeod"
BIN_TX="arkeo"
USER="$1"
SERVICE="$2"
BOND="$3"

PUBKEY_RAW=$($BIN keys show "$USER" -p --keyring-backend test | jq -r .key)
PUBKEY=$($BIN debug pubkey-raw "$PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

$BIN tx $BIN_TX bond-provider -y -b block --from "$USER" --keyring-backend test -- "$PUBKEY" "$SERVICE" "$BOND"
