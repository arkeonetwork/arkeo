#!/bin/bash

set -o pipefail
set -ex

CHAIN_ID="arkeo"
RPC="${RPC:=seed31.innovationtheory.com:26657}"
SEED="${SEED:=seed31.innovationtheory.com:26656}"
PEER_ID=$(curl -sL http://$RPC/status | jq -r '.result.node_info.id')
GENESIS="${GENESIS:=$RPC}"

if [ ! -f ~/.arkeo/config/genesis.json ]; then
	echo "setting validator node"

	arkeod init local --chain-id "$CHAIN_ID"
   	arkeod config set client keyring-backend test
   	arkeod config set client chain-id arkeo-testnet-3
	

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

	sed -i -e  "s/^pruning *=.*/pruning = \"custom\"/" $HOME/.arkeo/config/app.toml
   	sed -i -e  "s/^pruning-keep-recent *=.*/pruning-keep-recent = \"100\"/" $HOME/.arkeo/config/app.toml
   	sed -i -e  "s/^pruning-interval *=.*/pruning-interval = \"50\"/" $HOME/.arkeo/config/app.toml
   	sed -i -e  "s|^minimum-gas-prices =.*|minimum-gas-prices = \"0.001uarkeo\"|g" $HOME/.arkeo/config/app.toml
   	sed -i -e  "s|swagger =.*| swagger = true|g" $HOME/.arkeo/config/app.toml
	sed -i -e  "s/seeds = \"\"/seeds = \"$PEER_ID@$SEED\"/g" ~/.arkeo/config/config.toml

	# TODO: create this one as a validator
	# arkeod tx staking create-validator --chain-id arkeo --commission-rate 0.05 --commission-max-rate 0.2 --commission-max-change-rate 0.1 --min-self-delegation "1" --amount <staking amount>uarkeo --pubkey $(arkeod tendermint show-validator) --moniker "<yourvalidator-name>" --from <your-wallet-name> --fees="5000uarkeo" --yes
fi

arkeod start
