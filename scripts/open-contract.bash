#!/bin/bash
set -ex
if [ -z "$1" ]; then
	echo "No user supplied"
	exit 1
fi

if [ -z "$2" ]; then
	echo "No provider supplied"
	exit 1
fi

if [ -z "$3" ]; then
	echo "No service supplied"
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

BIN="arkeod"
BIN_TX="arkeo"
USER="$1"
PROVIDER="$2"
SERVICE="$3"
CTYPE="$4"
DURATION="$5"
RATE="15"
QUERY_PER_MINUTE=10
SETTLEMENT_DURATION=10

PUBKEY_RAW=$($BIN keys show "$PROVIDER" -p --keyring-backend test | jq -r .key)
PUBKEY=$($BIN debug pubkey-raw "$PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

CLIENT_PUBKEY_RAW=$($BIN keys show "$USER" -p --keyring-backend test | jq -r .key)
CLIENT_PUBKEY=$($BIN debug pubkey-raw "$CLIENT_PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

# nolint
DEPOSIT=$((RATE * DURATION))

$BIN tx $BIN_TX open-contract -y --from "$USER" --keyring-backend test -- "$PUBKEY" "$SERVICE" "$CLIENT_PUBKEY" "$CTYPE" "$DEPOSIT" "$DURATION" "$RATE""uarkeo" $QUERY_PER_MINUTE $SETTLEMENT_DURATION
