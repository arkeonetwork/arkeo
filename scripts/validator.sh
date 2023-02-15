#!/bin/sh

set -o pipefail
set -ex

PEER="${PEER:=none}" # the hostname of a seed node set as tendermint persistent peer
CHAIN_ID="arkeo"
TOKEN="uarkeo"
PORT_RPC=${PORT_RPC:=26657}
PORT_P2P=${PORT_P2P:=26656}

if [ ! -f ~/.arkeo/config/genesis.json ]; then
	echo "setting validator node"

	arkeod init local --chain-id "$CHAIN_ID"

	rm -rf ~/.arkeo/config/genesis.json

	if [ "$PEER" = "none" ]; then
		echo "Missing PEER"
		exit 1
	fi

	# wait for peer
	until curl -s "$PEER:$PORT_RPC" 1>/dev/null 2>&1; do
		echo "Waiting for peer: $PEER:$PORT_RPC"
		sleep 3
	done

	# fetch genesis file from seed node
	curl "$PEER:$PORT_RPC/genesis" | jq '.result.genesis' >~/.arkeo/config/genesis.json

	# fetch node id
	SEED_ID=$(curl -s "$PEER:$PORT_RPC/status" | jq -r .result.node_info.id)
	SEEDS="$SEED_ID@$PEER:$PORT_P2P"

	sed -i 's/enable = false/enable = true/g' ~/.arkeo/config/app.toml
	sed -i "s/seeds = \"\"/seeds = \"$SEEDS\"/g" ~/.arkeo/config/config.toml
	# TODO: create this one as a validator
	# arkeod tx staking create-validator --amount=100000000000uarkeo --pubkey=$(arkeod tendermint show-validator) --moniker="validator 1" --from=bob --keyring-backend test --commission-rate="0.10" --commission-max-rate="0.20" --commission-max-change-rate="0.01" --min-self-delegation="1"
fi

arkeod start
