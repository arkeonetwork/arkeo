#!/bin/bash

if [ -z "$1" ]; then
    echo "No user supplied"
    exit 1
fi

BIN="mercuryd"
BIN_TX="mercury"
TOKEN="token"
USER="$1"

ADDRESS=$($BIN keys show $USER -a)

$BIN query bank balances --denom $TOKEN -o json -- $ADDRESS
