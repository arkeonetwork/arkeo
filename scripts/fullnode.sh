#!/bin/bash

set -o pipefail

export RPC="${RPC:=seed.arkeo.network:26657}"
export SEED="${SEED:=seed.arkeo.network:26656}"

if [ ! -f ~/.arkeo/config/genesis.json ]; then
	CHAIN_ID=$(curl -s http://$RPC/status | jq -r '.result.node_info.network')
	arkeod init local --chain-id "$CHAIN_ID"

	curl -s http://$RPC/genesis | jq -r '.result.genesis' >~/.arkeo/config/genesis.json
fi

PEER_ID=$(curl -s http://$RPC/status | jq -r '.result.node_info.id')
sed -i "s/seeds = \"\"/seeds = \"$PEER_ID@$SEED\"/g" ~/.arkeo/config/config.toml

exec arkeod start
