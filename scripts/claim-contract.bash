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
    echo "No nonce supplied"
    exit 1
fi

if ! command -v signhere &> /dev/null
then
    echo "'signhere' could not be found. Run `make tools` and try again"
    exit 1
fi

BIN="mercuryd"
BIN_TX="mercury"
USER="$1"
PROVIDER="$2"
CHAIN="$3"
NONCE="$4"

PUBKEY_RAW=$($BIN keys show $PROVIDER -p | jq -r .key)
PUBKEY=$($BIN debug pubkey-raw $PUBKEY_RAW | grep "Bech32 Acc" | awk '{ print $NF }')

CLIENT_PUBKEY_RAW=$($BIN keys show $USER -p | jq -r .key)
CLIENT_PUBKEY=$($BIN debug pubkey-raw $CLIENT_PUBKEY_RAW | grep "Bech32 Acc" | awk '{ print $NF }')

HEIGHT=$(curl -s localhost:1317/mercury/contract/$PUBKEY/$CHAIN/$CLIENT_PUBKEY | jq -r .contract.height)

if [ -z "$HEIGHT" ]; then
    echo "No open contract to claim"
    exit 1
fi
if [ "$HEIGHT" == "null" ]; then
    echo "No open contract to claim"
    exit 1
fi

SIGNATURE=$(signhere -u $USER -m "$PUBKEY:$CHAIN:$CLIENT_PUBKEY:$HEIGHT:$NONCE")

$BIN tx $BIN_TX claim-contract-income -y --from $USER --gas auto -- $PUBKEY $CHAIN $CLIENT_PUBKEY $NONCE $HEIGHT $SIGNATURE
