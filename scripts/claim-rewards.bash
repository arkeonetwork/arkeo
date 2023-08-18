#!/bin/bash

if [ -z "$1" ]; then
	echo "No host supplied"
	exit 1
fi

if [ -z "$2" ]; then
	echo "No user supplied"
	exit 1
fi

# should be an ip address or domain name
HOST=$1

# should be the key name (ie arkeod keys list --keyring-backend test)
USER=$2

BIN="arkeod"
BIN_TX="arkeo"

raw_claims=$(curl -sL "$HOST":3636/open-claims | jq '.[] | select(.claimed == false) | @json')
for raw_claim in $raw_claims; do

	id=$(echo "$raw_claim" | jq -r '.contract_id')
	nonce=$(echo "$raw_claims" | jq -r '.nonce')
	signature=$(echo "$raw_claims" | jq -r '.signature')

	$BIN tx $BIN_TX claim-contract-income -y -b block --from "$USER" --keyring-backend test -- "$id" "$nonce" "$signature"
done
