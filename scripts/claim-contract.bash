#!/bin/bash

if [ -z "$1" ]; then
	echo "No user supplied"
	exit 1
fi

if [ -z "$2" ]; then
	echo "No contract id supplied"
	exit 1
fi

if [ -z "$3" ]; then
	echo "No nonce supplied"
	exit 1
fi

if ! command -v signhere &>/dev/null; then
	echo "'signhere' could not be found. Run $(make tools) and try again"
	exit 1
fi

BIN="arkeod"
BIN_TX="arkeo"
USER="$1"
CONTRACT_ID="$2"
NONCE="$3"

CLIENT_PUBKEY_RAW=$($BIN keys show "$USER" -p --keyring-backend test | jq -r .key)
CLIENT_PUBKEY=$($BIN debug pubkey-raw "$CLIENT_PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

HEIGHT=$(curl -s localhost:1317/arkeo/contract/"$CONTRACT_ID" | jq -r .contract.height)

if [ -z "$HEIGHT" ]; then
	echo "No open contract to claim"
	exit 1
fi
if [ "$HEIGHT" == "null" ]; then
	echo "No open contract to claim"
	exit 1
fi

SIGNATURE=$(signhere -u "$USER" -m "$CONTRACT_ID:$CLIENT_PUBKEY:$HEIGHT:$NONCE")

$BIN tx $BIN_TX claim-contract-income -y --from "$USER" --keyring-backend test -- "$CONTRACT_ID" "$CLIENT_PUBKEY" "$NONCE" "$HEIGHT" "$SIGNATURE"
