# [Round 1] Testing Docs

This phase of testing is about getting the complete full stack of a working environment, run by a single dev, while the team can run through a series of testing.

## Step 1: get the cli installed
Either clone and build arkeo and its tools from source as outlined in [readme.md](../readme.md),
or download a binary for your operating system below:

- [macOS (x86-64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/darwin_amd64/arkeod)
  - sha256 `427f3edfd0d7d58719f8a33d65826d39fa45a9fa2fa4e5e70835d8b4117b8ef0`
- [macOS (arm64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/darwin_arm64/arkeod)
  - sha256 `3d55c33393aa744fbc619b70fe802413394503e6c941023316f99137d8944792`
- [linux (x86-64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/linux_amd64/arkeod)
  - sha256 `12a66411c342c0874778066a9548ff220850deb4265c79ccfa98e508d939d80d`
- [linux (arm64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/linux_arm64/arkeod)
  - sha256 `69acb1916a5715fbf00eb733252999665d4f71a8a35c07e2c8e360cb7d78caad`

after downloading the executable, verify the integrity of the downloaded artifact:
```bash
$ cd /path/to/downloads
# replace the sha256 sum we echo below with the appropriate sum for your os listed above
$ echo '12a66411c342c0874778066a9548ff220850deb4265c79ccfa98e508d939d80d arkeod' | sha256sum -c -       
arkeod: OK
```
note the output "`arkeod: OK`"

now move the `arkeod` file to a directory that's on your PATH, or add the containing directory to your PATH:
```bash
$ export PATH=$PATH:/path/to/directory_containing
```
/path/to/directory is the directory you placed the arkeod binary you downloaded in. to make this permanent,
add the `export PATH=...` statement above to your shell initialization scripts (.bashrc, .zshrc, etc.)

verify installation by executing the arkeod command:
```
$ arkeod version
0.0.1
```

then update your client config as follows:
```bash
$ arkeod config chain-id arkeo
$ arkeod config node tcp://testnet-seed.arkeo.shapeshift.com:26657
```
you can verify the configuration applied successfully:
```bash
$ arkeod config
{
  "chain-id": "arkeo",
  "keyring-backend": "os",
  "output": "text",
  "node": "tcp://testnet-seed.arkeo.shapeshift.com:26657",
  "broadcast-mode": "sync"
}
```
and query the current block height:
```
$ arkeod query block | jq '.block.header.height'
"157326"
```
## Step 2: Setup a wallet
Create a wallet using the cli and find your address and pubkey.

```bash
$ arkeod keys add <user> --keyring-backend file
```
__Note__: optionally recover from an existing bip39 mnemonic by adding `--recover` to the `keys add` command above

You should see your address from the output of this command. It will start with `arkeo1`

To find your pubkey, take the `key` of the pubkey in the output of our last command

```bash
$ arkeod debug pubkey-raw <raw pubkey>
```

This should output your pubkey. It will start with `arkeopub1`.

## Step 2: Get funded

Reach out to Odin with your address and request for funds sent to you.

## Step 3: Open a Contract

Open both a subscription and pay-as-you-go contract (not at the same time, as only one contract is allowed open at a time between a provider and client/delegate) with the sole data provider.

```bash
$ CTYPE=0 # contract type, 0 is subscription, 1 is pay-as-you-go
$ DEPOSIT=<amt> # amount of tokens you want to deposit. Subscriptions should make sense in that duration and rate equal deposit
$ DURATION=<blocks> # number of blocks to make a subscription. There are lower and higher limits to this number
$ RATE=<rate> # should equal the porvider's rate which you can lookup at (`curl http://seed.arkeo.network:3636/metadata.json | jq .`)
$ arkeod tx arkeo open-contract -y --from <user> --keyring-backend file --node "tcp://seed.arkeo.network:26657" -- arkeopub1addwnpepqtrc0rrpkwn2esula68zl3dvqqfxfjhr5dyfxy3uq97dssntrq8twhy9nvu btc-mainnet-fullnode "<your pubkey>" "$CTYPE" "$DEPOSIT" "$DURATION" $RATE
```

## Step 4: Make Requests

There is a simple cli tool that was written. To install it, use `make tools`. The tool, `curleo`, abstracts away all of the signing and authentication elements for the user to simplify making an authenticated request.

```bash
$ curleo -u bob -data '{"jsonrpc": "1.0", "id": "curltest", "method": "ping", "params": []}' -H "text/plain" http://seed.arkeo.network:3636/btc-mainnet-fullnode | jq
{
  "result": null,
  "error": null,
  "id": "curltest"
}
```

## Step 4: Claim Rewards for the Provider

On the provider’s behalf, you can claim the rewards for them (for testing purposes). To do use the following command. 

```bash
$ NONCE=<num> # the nonce represents the number of queries made between the client/provider and provider during this contract
$ HEIGHT=<height> # the block height the contract was open
$ SIGNATURE=$(signhere -u <user> -m "arkeopub1addwnpepqtrc0rrpkwn2esula68zl3dvqqfxfjhr5dyfxy3uq97dssntrq8twhy9nvu:btc-mainnet-fullnode:<your pubkey>:$HEIGHT:$NONCE") # signature
$ arkeod tx arkeo claim-contract-income -y --from <user> --keyring-backend file --node "tcp://seed.arkeo.network:26657" -- arkeopub1addwnpepqtrc0rrpkwn2esula68zl3dvqqfxfjhr5dyfxy3uq97dssntrq8twhy9nvu btc-mainnet-fullnode <your pubkey> "$NONCE" "$HEIGHT" "$SIGNATURE"
```

## Step 5: Close a Contract

If the contract is a subscription, it can be cancelled. Pay-as-you-go isn’t available to cancel as you can stop making requests as a form of cancelling (providers can cancel though). Closing a contract should also trigger a payout to the provider.

```bash
$ arkeod tx arkeo close-contract -y --from <user> --keyring-backend file --node "tcp://seed.arkeo.network:26657" -- arkeopub1addwnpepqtrc0rrpkwn2esula68zl3dvqqfxfjhr5dyfxy3uq97dssntrq8twhy9nvu btc-mainnet-fullnode "<your pubkey>"
```