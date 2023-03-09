# [Round 1] Testing Docs

This phase of testing is about getting the complete full stack of a working environment, run by a single dev, while the team can run through a series of testing.

## Step 1: get the cli installed
Either clone and build arkeo and its tools from source as outlined in [readme.md](../readme.md),
or download a binary for your operating system below:

- [macOS (x86-64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_darwin_amd64.tar.gz)
  - [sha256](sums/arkeo_darwin_amd64.sha256)
- [macOS (arm64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_darwin_arm64.tar.gz)
  - [sha256](sums/arkeo_darwin_arm64.sha256)
- [linux (x86-64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_linux_arm64.tar.gz)
  - [sha256](sums/arkeo_linux_arm64.sha256)
- [linux (arm64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_linux_arm64.tar.gz)
  - [sha256](sums/arkeo_linux_arm64.sha256)

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

optionally set the backend keyring. the default value "os" uses the operating system's keyring.
```bash
$ arkeod config keyring-backend <os|test|file>
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
$ arkeod query block | jq -r '.block.header.height'
205160
```
## Step 2: Setup a wallet
Create or Import a wallet using the cli. Replace `adam` below with a name of your choice. add `--recover` if you have a
mnemonic you'd like to use.

assign a name to the $ark_user variable. this will become the key/wallet/user name.
```bash
$ ark_user=adam
```

```bash
$ arkeod keys add $ark_user
```
__or__
```bash
$ arkeod keys add $ark_user --recover
```
This will output your arkeo address and pubkey along with the name you selected, as well as the mnemonic
if you opted to generate one. Save this somewhere safe.

```
- address: arkeo14q4qnm4tmkm9xuhjwu0vw0f8xy7ztexeesvflj
  name: adam
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"ApcnOwEDoO6smze46IPUgC/5bC8DohEpLJ9ZZnrKky0w"}'
  type: local
```

In order to interact with arkeo providers, you will need to have the `Acc` (account) pubkey encoded bech32 with the standard prefix. Execute the command below to obtain it.

```bash
$ arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r '.key') | grep '^Bech32 Acc: ' | awk '{ print $NF }'
arkeopub1addwnpepq2tjwwcpqwswatymx7uw3q75sqhljmp0qw3pz2fvnavkv7k2jvknq9k9lr0
```

Store your menemonic somewhere safe along with the pubkey (`arkeopub1addwnpepq...9lr0` above) and your address.

## Step 2: Get funded

Reach out to the arkeo development team with your address (starts with `arkeo1`) to request funds.

## Step 3: Open a Contract

Get a list of Online providers from the Directory Service:
```bash
$ curl -s http://directory.arkeo.shapeshift.com/provider/search/ | jq '.[]|select(.Status == "Online")|[{pubkey: .Pubkey, chain: .Chain, meta: .MetadataURI}]'
[
  {
    "pubkey": "arkeopub1addwnpepqdtyf722w22r8grkecpnzgwm6stm2x3yhphre2wwnwwxkpa9ym5fyfyxdum",
    "chain": "gaia-mainnet-rpc-archive",
    "meta": "http://testnet-sentinel.arkeo.shapeshift.com:3636/metadata.json"
  }
]
```

Choose the `gaia-mainnet-rpc-archive` provider:
```bash
$ ark_provider=arkeopub1addwnpepqdtyf722w22r8grkecpnzgwm6stm2x3yhphre2wwnwwxkpa9ym5fyfyxdum
$ ark_chain=gaia-mainnet-rpc-archive
```

Open a Subscription contract. This example opens a subscription contract for 20 blocks at a rate of 10 arkeo, depositing 200 to cover the subscription cost.
```bash
$ ark_contract_type=0
$ ark_deposit=200
$ ark_duration=20
$ ark_rate=10
$ ark_pubkey=`arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r '.key') | grep '^Bech32 Acc: ' | awk '{ print $NF }'`
$ arkeod tx arkeo open-contract --from $ark_user -- $ark_provider $ark_chain "$ark_pubkey" "$ark_contract_type" "$ark_deposit" "$ark_duration" $ark_rate
```

Open a Pay-As-You-Go contract. This example opens a subscription contract for 20 blocks at a rate of 20 arkeo, depositing 400 to cover the subscription cost.
```bash
$ ark_contract_type=1
$ ark_deposit=400
$ ark_duration=20
$ ark_rate=20
$ ark_pubkey=`arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r '.key') | grep '^Bech32 Acc: ' | awk '{ print $NF }'`
$ arkeod tx arkeo open-contract --from $ark_user -- $ark_provider $ark_chain "$ark_pubkey" "$ark_contract_type" "$ark_deposit" "$ark_duration" $ark_rate
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