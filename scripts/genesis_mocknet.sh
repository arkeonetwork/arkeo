#!/bin/bash

set -o pipefail
set -ex

CHAIN_ID="arkeo-main-v2"
STAKE="100000000000uarkeo"
TOKEN="uarkeo"
USER="genesis"
TOTAL_SUPPLY=100000000000 # Initial supply corresponding to the stake

add_module() {
	jq --arg ADDRESS "$1" --arg ASSET "$2" --arg AMOUNT "$3" --arg NAME "$4" '.app_state.auth.accounts += [{
        "@type": "/cosmos.auth.v1beta1.ModuleAccount",
        "base_account": {
          "address": $ADDRESS,
          "pub_key": null,
          "sequence": "0"
        },
        "name": $NAME,
        "permissions": []
  }]' <~/.arkeo/config/genesis.json >/tmp/genesis.json
	mv /tmp/genesis.json ~/.arkeo/config/genesis.json

	jq --arg ADDRESS "$1" --arg ASSET "$2" --arg AMOUNT "$3" '.app_state.bank.balances += [{
        "address": $ADDRESS,
        "coins": [ { "denom": $ASSET, "amount": $AMOUNT } ],
    }]' <~/.arkeo/config/genesis.json >/tmp/genesis.json
	mv /tmp/genesis.json ~/.arkeo/config/genesis.json

	TOTAL_SUPPLY=$((TOTAL_SUPPLY + $3))
	echo "Total supply after adding module: $TOTAL_SUPPLY"
}

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

	TOTAL_SUPPLY=$((TOTAL_SUPPLY + $3))
	echo "Total supply after adding account: $TOTAL_SUPPLY"
}

add_claim_records() {
	jq --arg CHAIN "$1" --arg ADDRESS "$2" --arg AMOUNTCLAIM "$3" --arg AMOUNTVOTE "$4" --arg AMOUNTDELEGATE "$5" --arg ISTRANSFERABLE "$6" '.app_state.claimarkeo.claim_records += [{
        "chain": $CHAIN,
		"address": $ADDRESS,
        "amount_claim": { "denom": "uarkeo", "amount": $AMOUNTCLAIM },
        "amount_vote": { "denom": "uarkeo", "amount": $AMOUNTVOTE },
        "amount_delegate": { "denom": "uarkeo", "amount": $AMOUNTDELEGATE },
        "is_transferable": true,
    }]' <~/.arkeo/config/genesis.json >/tmp/genesis.json
	mv /tmp/genesis.json ~/.arkeo/config/genesis.json
}

set_fee_pool() {
	local denom="$1"
	local amount="$2"

	jq --arg DENOM "$denom" --arg AMOUNT "$amount" '.app_state.distribution.fee_pool.community_pool = [{
        "denom": $DENOM,
        "amount": $AMOUNT
    }]' <~/.arkeo/config/genesis.json >/tmp/genesis.json
	mv /tmp/genesis.json ~/.arkeo/config/genesis.json
}

disable_mint_params() {
	jq '.app_state.mint.minter.inflation = "0.000000000000000000" |
        .app_state.mint.minter.annual_provisions = "0.000000000000000000" |
        .app_state.mint.params.inflation_rate_change = "0.000000000000000000" |
        .app_state.mint.params.inflation_max = "0.000000000000000000" |
        .app_state.mint.params.inflation_min = "0.000000000000000000"' \
		<~/.arkeo/config/genesis.json >/tmp/genesis.json
	mv /tmp/genesis.json ~/.arkeo/config/genesis.json
}

if [ ! -f ~/.arkeo/config/priv_validator_key.json ]; then
	# remove the original generate genesis file, as below will init chain again
	rm -rf ~/.arkeo/config/genesis.json
fi

if [ ! -f ~/.arkeo/config/genesis.json ]; then

  arkeod init local --default-denom $TOKEN --chain-id "$CHAIN_ID"
  arkeod keys add $USER --keyring-backend test
  arkeod genesis add-genesis-account $USER $STAKE --keyring-backend test
  arkeod keys list --keyring-backend test
  arkeod genesis gentx $USER $STAKE --chain-id $CHAIN_ID --keyring-backend test
  arkeod genesis collect-gentxs

  # arkeod keys add faucet --keyring-backend test
  # FAUCET=$(arkeod keys show faucet -a --keyring-backend test)
  # add_account "$FAUCET" $TOKEN 2900000000000000 # faucet, 29m

  disable_mint_params

  # Reserve, 20%
  add_module arkeo1d0m97ywk2y4vq58ud6q5e0r3q9khj9e3nksryg $TOKEN 2420000000000000 'arkeo-reserve'
  # Claims/Airdrop, 20%
  add_module arkeo14tmx70mvve3u7hfmd45vle49kvylk6s2p6xuz8 $TOKEN 3025000000000000 'claimarkeo'
  # Distribution, 45%
  add_module arkeo1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8yl3e56 $TOKEN 5445000000000000 'distribution'
  # this is to handle the balance set to distribution as a community pool
  set_fee_pool $TOKEN 5445000000000000 # Community ' 45% '

  # Add Foundational Accounts
  add_account "arkeo1n2g60ykpmurqwsxce44e3l3dt6lsm2mv7u37a5" $TOKEN 1209600000000000 #' 10% '
  # validator 1 account
  add_account "arkeo1w2ln0prejgrztmf9w23e0rsnlks7djneh5te7p" $TOKEN 100000000000
  # validator 2 account
  add_account "arkeo1gcrgw7hpqn4aqk8uc3x9uqzg8mlqv0wuu4setj" $TOKEN 100000000000
  # validator 3 account
  add_account "arkeo1wtws43medp3ay8gm8nn4uyx3efdsqsj2yg6n2k" $TOKEN 100000000000

  # enable CORs on testnet/localnet
  # sed -i '' 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' ~/.arkeo/config/app.toml
  # sed -i '' 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/g' ~/.arkeo/config/config.toml

	# else
	#	 source ./add-claim-records.sh
	# fi

	sed -i '' 's/"stake"/"uarkeo"/g' ~/.arkeo/config/genesis.json
	sed -i '' '/"duration_until_decay"\|"duration_of_decay"/s/"3600s"/"7884000s"/' ~/.arkeo/config/genesis.json
	sed -i '' 's/enable = false/enable = true/g' ~/.arkeo/config/app.toml
	sed -i '' 's/127.0.0.1:26657/0.0.0.0:26657/g' ~/.arkeo/config/config.toml
	sed -i '' 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/g' ~/.arkeo/config/app.toml

	# Update the supply field in genesis.json using jq
	jq --arg DENOM "$TOKEN" --arg AMOUNT "$TOTAL_SUPPLY" '.app_state.bank.supply = [{"denom": $DENOM, "amount": $AMOUNT}]' <~/.arkeo/config/genesis.json >/tmp/genesis.json
	mv /tmp/genesis.json ~/.arkeo/config/genesis.json

	set -e
	arkeod validate-genesis --trace

fi

arkeod start --pruning nothing --minimum-gas-prices 0uarkeo
