#!/bin/bash

set -o pipefail
set -ex

CHAIN_ID="arkeo"
RPC="${RPC:=seed.arkeo.network:26657}"
SEED="${SEED:=seed.arkeo.network:26656}"
PEER_ID=$(curl -sL http://$RPC/status | jq -r '.result.node_info.id')
GENESIS="${GENESIS:=$RPC}"

if [ ! -f ~/.arkeo/config/genesis.json ]; then
	echo "setting validator node"

	arkeod init local --chain-id "$CHAIN_ID"

	rm -rf ~/.arkeo/config/genesis.json

	if [ "$RPC" = "none" ]; then
		echo "Missing PEER"
		exit 1
	fi

	# wait for peer
	until curl -s "$RPC" 1>/dev/null 2>&1; do
		echo "Waiting for peer: $RPC"
		sleep 3
	done

	# fetch genesis file from seed node
	curl -sL "$GENESIS/genesis" | jq '.result.genesis' >~/.arkeo/config/genesis.json

	# fetch node id
	SEED_ID=$(curl -sL "$RPC/status" | jq -r .result.node_info.id)
	SEEDS="$SEED_ID@$SEED"

	sed -i 's/enable = false/enable = true/g' ~/.arkeo/config/app.toml
	sed -i "s/seeds = \"\"/seeds = \"$PEER_ID@$SEED\"/g" ~/.arkeo/config/config.toml
	# TODO: create this one as a validator
	# arkeod tx staking create-validator --amount=100000000000uarkeo --pubkey=$(arkeod tendermint show-validator) --moniker="validator 1" --from=bob --keyring-backend test --commission-rate="0.10" --commission-max-rate="0.20" --commission-max-change-rate="0.01" --min-self-delegation="1"
fi

arkeod start
