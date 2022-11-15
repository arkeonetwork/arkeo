#!/bin/sh

set -o pipefail
set -x

CHAIN_ID="arkeo"
STAKE="100000000000ukeo"
TOKEN="ukeo"
USER="ark"

add_account() {
  jq --arg ADDRESS "$1" --arg ASSET "$2" --arg AMOUNT "$3" '.app_state.auth.accounts += [{
        "@type": "/cosmos.auth.v1beta1.BaseAccount",
        "address": $ADDRESS,
        "pub_key": null,
        "account_number": "0",
        "sequence": "0"
    }]' <~/.arkeo/config/genesis.json >/tmp/genesis.json
  mv /tmp/genesis.json ~/.arkeo/config/genesis.json

  jq --arg ADDRESS "$1" --arg ASSET "$2" --arg AMOUNT "$3" '.app_state.bank.balances += [{
        "address": $ADDRESS,
        "coins": [ { "denom": $ASSET, "amount": $AMOUNT } ],
    }]' <~/.arkeo/config/genesis.json >/tmp/genesis.json
  mv /tmp/genesis.json ~/.arkeo/config/genesis.json
}

if [ ! -f ~/.arkeo/config/priv_validator_key.json ]; then
	# remove the original generate genesis file, as below will init chain again
	rm -rf ~/.arkeo/config/genesis.json
fi

if [ ! -f ~/.arkeo/config/genesis.json ]; then
	arkeod init local --staking-bond-denom $TOKEN --chain-id "$CHAIN_ID"
	arkeod keys add $USER --keyring-backend test
	arkeod add-genesis-account $USER $STAKE --keyring-backend test
	arkeod keys list --keyring-backend test
	arkeod gentx $USER $STAKE --chain-id $CHAIN_ID --keyring-backend test
	arkeod collect-gentxs

    if [ "$NET" = "mocknet" ] || [ "$NET" = "testnet" ]; then
        add_account rko1dheycdevq39qlkxs2a6wuuzyn4aqxhvee2kjas ukeo 10000000000000000 # reserve, 100m

        arkeod keys add faucet --keyring-backend test
        FAUCET=$(arkeod keys show faucet -a)
        add_account $FAUCET ukeo 10000000000000000 # faucet, 100m
    fi

	sed -i 's/"stake"/"ukeo"/g' ~/.arkeo/config/genesis.json
	sed -i 's/enable = false/enable = true/g' ~/.arkeo/config/app.toml
	sed -i 's/127.0.0.1:26657/0.0.0.0:26657/g' ~/.arkeo/config/config.toml

	arkeod validate-genesis --trace
fi

arkeod start
