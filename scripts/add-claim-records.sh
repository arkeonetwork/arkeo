#!/bin/bash

# Input files and genesis file
CHAIN1_CSV="./data/arkeo_airdrop_bech32.csv" # Replace with the actual filename for the first chain
CHAIN2_CSV="./data/combined_eth_airdrop.csv" # Replace with the actual filename for the second chain
GENESIS_FILE="~/.arkeo/config/genesis.json"
TEMP_FILE="/tmp/genesis_temp.json"

# Check if CSV files and genesis.json exist
if [ ! -f "$CHAIN1_CSV" ] || [ ! -f "$CHAIN2_CSV" ]; then
    echo "CSV file for one of the chains not found!"
    exit 1
fi

# if [ ! -f "$GENESIS_FILE" ]; then
#     echo "Genesis file not found!"
#     exit 1
# fi

# Function to process a CSV file and add entries to the genesis file
add_claim_records() {
    local csv_file="$1"
    local chain_name="$2"

    # Start building the new claim records JSON array
    records="["

    # Loop through each line in the CSV file (skip the header)
    while IFS=, read -r address amount; do
        # Multiply the amount by 10^9
        amount_claim=$(echo "scale=0; $amount * 100000000 / 3" | bc)
        amount_vote=$(echo "scale=0; $amount * 100000000 / 3" | bc)
        amount_delegate=$(echo "scale=0; $amount * 100000000 / 3" | bc)

        # Append the new record to the JSON array
        records+=$(jq -n --arg chain "$chain_name" --arg address "$address" --arg amount_claim "$amount_claim" --arg amount_vote "$amount_vote" --arg amount_delegate "$amount_delegate" \
            '{chain: $chain, address: $address, amount_claim: {denom: "uarkeo", amount: $amount_claim}, amount_vote: {denom: "uarkeo", amount: $amount_vote}, amount_delegate: {denom: "uarkeo", amount: $amount_delegate}, is_transferable: true}')
        records+=","
    done < <(tail -n +2 "$csv_file")

    # Remove the trailing comma and close the JSON array
    records="${records%,}]"

    echo "Final JSON Records Array: $records"

    # Add the batch of records to the genesis.json in one operation
    jq --argjson new_records "$records" '.app_state.claimarkeo.claim_records += $new_records' <~/.arkeo/config/genesis.json >/tmp/genesis.json
    mv /tmp/genesis.json ~/.arkeo/config/genesis.json
}

# Process CSV files for each chain
add_claim_records "$CHAIN1_CSV" "ARKEO"
add_claim_records "$CHAIN2_CSV" "ETHEREUM"

echo "Updated genesis.json with new claim records for both chains."
