#!/bin/bash

USER="alice"

echo "Waiting for RPC..."
until curl -s "$EVENT_STREAM_HOST/status" >/dev/null; do
	# echo "Rest server is unavailable - sleeping"
	sleep 5
done

if [ "$NET" = "mocknet" ] || [ "$NET" = "testnet" ]; then
	PUBKEY_RAW=$(arkeod keys show "$USER" -p --keyring-backend test | jq -r .key)
	PUBKEY=$(arkeod debug pubkey-raw "$PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')
fi

PROVIDER_PUBKEY=$PUBKEY exec sentinel
