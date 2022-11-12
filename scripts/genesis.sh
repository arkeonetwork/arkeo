#!/bin/sh

set -o pipefail

CHAIN_ID="arkeo"
STAKE="100000000000stake"
USER="ark"

if [ ! -f ~/.arkeo/config/priv_validator_key.json ]; then
	# remove the original generate genesis file, as below will init chain again
	rm -rf ~/.arkeo/config/genesis.json
fi

if [ ! -f ~/.arkeo/config/genesis.json ]; then
	arkeod init local --chain-id "$CHAIN_ID"
	arkeod keys add $USER --keyring-backend test
	arkeod add-genesis-account $USER $STAKE --keyring-backend test
	arkeod keys list --keyring-backend test
	arkeod gentx $USER $STAKE --chain-id $CHAIN_ID --keyring-backend test
	arkeod collect-gentxs

	arkeod validate-genesis --trace
	sed -i 's/enable = false/enable = true/g' ~/.arkeo/config/app.toml
fi

arkeod start
