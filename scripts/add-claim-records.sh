#!/bin/bash

# Input files and genesis file
CHAIN1_CSV="./data/arkeo_airdrop.csv" 
CHAIN2_CSV="./data/eth_airdrop.csv"
GENESIS_FILE="~/.arkeo/config/genesis.json"
TEMP_FILE="/tmp/genesis_temp.json"

# Check if CSV files and genesis.json exist
if [ ! -f "$CHAIN1_CSV" ] || [ ! -f "$CHAIN2_CSV" ]; then
    echo "CSV file for one of the chains not found!"
    exit 1
fi

# Function to process a CSV file and add entries to the genesis file
add_claim_records() {
    local csv_file="$1"
    local chain_name="$2"

    # Start building the new claim records JSON array
    records="["

    # Loop through each line in the CSV file (skip the header)
    while IFS=, read -r address amount; do
        amount_claim=$(echo "scale=0; $amount * 100000000 / 3" | bc)
        amount_vote=$(echo "scale=0; $amount * 100000000 / 3" | bc)
        amount_delegate=$(echo "scale=0; $amount * 100000000 / 3" | bc)

        # Append the new record to the JSON array
        records+=$(jq -n --arg chain "$chain_name" --arg address "$address" --arg amount_claim "$amount_claim" --arg amount_vote "$amount_vote" --arg amount_delegate "$amount_delegate" \
            '{chain: $chain, address: $address, amount_claim: {denom: "uarkeo", amount: $amount_claim}, amount_vote: {denom: "uarkeo", amount: $amount_vote}, amount_delegate: {denom: "uarkeo", amount: $amount_delegate}, is_transferable: true}')
        records+=","
    done < <(tail -n +2 "$csv_file")

    records="${records%,}]"

    echo "Final JSON Records Array: $records"

    jq --argjson new_records "$records" '.app_state.claimarkeo.claim_records += $new_records' <~/.arkeo/config/genesis.json >/tmp/genesis.json
    mv /tmp/genesis.json ~/.arkeo/config/genesis.json
}

add_claim_records "$CHAIN1_CSV" "ARKEO"
add_claim_records "$CHAIN2_CSV" "ETHEREUM"

echo "Updated genesis.json with new claim records for arkeo and ethereum."
