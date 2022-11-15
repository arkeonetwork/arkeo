#!/bin/sh

USER="ark"

echo "Waiting for RPC..."
until curl -s "$EVENT_STREAM_HOST/status" >/dev/null; do
	# echo "Rest server is unavailable - sleeping"
	sleep 5
done

if [ "$NET" = "mocknet" ] || [ "$NET" = "testnet" ]; then
	while true; do
		if arkeod keys show "$USER" -p --keyring-backend test; then
			PUBKEY_RAW=$(arkeod keys show "$USER" -p --keyring-backend test | jq -r .key)
			if PUBKEY=$(arkeod debug pubkey-raw "$PUBKEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }'); then
				break
			fi
		fi
		sleep 3
	done
fi

PROVIDER_PUBKEY=$PUBKEY sentinel
