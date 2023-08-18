#!/bin/bash

if [ -z "$1" ]; then
	echo "No user supplied"
	exit 1
fi

if [ -z "$2" ]; then
	echo "No contract id supplied"
	exit 1
fi

BIN="arkeod"
BIN_TX="arkeo"
USER="$1"
ID="$2"

$BIN tx $BIN_TX close-contract -y -b block --from "$USER" --keyring-backend test -- "$ID"
