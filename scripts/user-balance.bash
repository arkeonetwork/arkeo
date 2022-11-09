#!/bin/bash

if [ -z "$1" ]; then
	echo "No user supplied"
	exit 1
fi

BIN="arkeod"
TOKEN="keo"
USER="$1"

ADDRESS=$($BIN keys show "$USER" -a)

$BIN query bank balances --denom $TOKEN -o json -- "$ADDRESS"
