# [Round 1] Testing Docs

This phase of testing is about getting the complete full stack of a working environment, run by a single dev, while the team can run through a series of testing.

## Step 1: get the cli installed
Either clone and build arkeo and its tools from source as outlined in [readme.md](../readme.md),
or download a binary for your operating system below:

- [macOS (x86-64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_darwin_amd64.tar.gz)
  - [sha256](sums/arkeo_darwin_amd64.sha256?raw=true)
- [macOS (arm64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_darwin_arm64.tar.gz)
  - [sha256](sums/arkeo_darwin_arm64.sha256?raw=true)
- [linux (x86-64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_linux_arm64.tar.gz)
  - [sha256](sums/arkeo_linux_arm64.sha256?raw=true)
- [linux (arm64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_linux_arm64.tar.gz)
  - [sha256](sums/arkeo_linux_arm64.sha256?raw=true)

after downloading the executable and corresponding checksum file (right click->save as/link/target), verify the integrity of the downloaded artifact:
```bash
# using macOS arm64 as the example, substitute your platform as needed. the below will extract arkeod, curleo, and signhere
# to /usr/local/bin. replace "-C /usr/local/bin" with a different path if desired.
cd /path/to/downloads
sha256sum -c arkeo_darwin_arm64.sha256
arkeo_darwin_arm64.tar.gz: OK
tar -C /usr/local/bin -zxvf arkeo_darwin_arm64.tar.gz # add sudo if needed
curleo
arkeod
signhere
```
note the output "`arkeod: OK`"

if the path you extracted the binaries to is not on your PATH (/usr/local/bin generally is on *nix), add the directory to your PATH:
```bash
export PATH=$PATH:/path/to/directory_containing
```
/path/to/directory is the directory you extracted the binaries you downloaded to with the tar command above. to make this permanent,
add the prior `export PATH=...` statement to your shell initialization scripts (.bashrc, .zshrc, etc.)

verify installation by executing the arkeod command:
```
arkeod version
0.0.1
```

then update your client config as follows:
```bash
arkeod config chain-id arkeo
arkeod config node tcp://testnet-seed.arkeo.shapeshift.com:26657
```

optionally set the backend keyring. the default value "os" uses the operating system's keyring.
```bash
arkeod config keyring-backend <os|test|file>
```

you can verify the configuration applied successfully:
```bash
arkeod config
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
arkeod query block | jq -r '.block.header.height'
205160
```
## Step 2: Setup a wallet
Create or Import a wallet using the cli. Replace `adam` below with a name of your choice. add `--recover` if you have a
mnemonic you'd like to use.

assign a name to the $ark_user variable. this will become the key/wallet/user name.
```bash
ark_user=adam
```

```bash
arkeod keys add $ark_user
```
__or__
```bash
arkeod keys add $ark_user --recover
```
This will output your arkeo address and pubkey along with the name you selected, as well as the mnemonic
if you opted to generate one. Save this somewhere safe.

```
- address: tarkeo14q4qnm4tmkm9xuhjwu0vw0f8xy7ztexek44nw3
  name: adam
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"ApcnOwEDoO6smze46IPUgC/5bC8DohEpLJ9ZZnrKky0w"}'
  type: local
```

In order to interact with arkeo providers, you will need to have the `Acc` (account) pubkey encoded bech32 with the standard prefix. Execute the command below to obtain it.

```bash
arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r '.key') | grep '^Bech32 Acc: ' | awk '{ print $NF }'
tarkeopub1addwnpepq2tjwwcpqwswatymx7uw3q75sqhljmp0qw3pz2fvnavkv7k2jvknqk25q7j
```

Store your menemonic somewhere safe along with the pubkey (`arkeopub1addwnpepq...9lr0` above) and your address.

## Step 2: Get funded

Reach out to the arkeo development team with your address (starts with `arkeo1`) to request funds.

## Step 3: Open a Contract

Get a list of Online providers from the Directory Service:
```bash
curl -s http://directory.arkeo.shapeshift.com/provider/search/ | jq '.[]|select(.Status == "ONLINE")|[{pubkey: .Pubkey, chain: .Chain, meta: .MetadataURI}]'
[
  {
    "pubkey": "tarkeopub1addwnpepq0h7hn9jzhkfwkxgp6kl3ljtjxfvz48emzdrrt5epzjrumpx9kz3w9mjsq9",
    "chain": "gaia-mainnet-rpc-archive",
    "meta": "http://testnet-sentinel.arkeo.shapeshift.com:3636/metadata.json"
  }
]
```

Choose the `gaia-mainnet-rpc-archive` provider:
```bash
ark_provider=tarkeopub1addwnpepq0h7hn9jzhkfwkxgp6kl3ljtjxfvz48emzdrrt5epzjrumpx9kz3w9mjsq9
ark_chain=gaia-mainnet-rpc-archive
```

Open a Subscription contract. This example opens a subscription contract for 20 blocks at a rate of 10 arkeo, depositing 200 to cover the subscription cost.
```bash
ark_contract_type=0
ark_deposit=200
ark_duration=20
ark_settle_duration=10
ark_rate=10
ark_pubkey=`arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r '.key') | grep '^Bech32 Acc: ' | awk '{ print $NF }'`
arkeod tx arkeo open-contract --from $ark_user -- $ark_provider $ark_chain "$ark_pubkey" "$ark_contract_type" "$ark_deposit" "$ark_duration" $ark_rate "$ark_settle_duration"
```

Open a Pay-As-You-Go contract. This example opens a subscription contract for 20 blocks at a rate of 20 arkeo, depositing 400 to cover the subscription cost.
```bash
ark_contract_type=1
ark_deposit=400
ark_duration=20
ark_settle_duration=10
ark_rate=20
ark_pubkey=`arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r '.key') | grep '^Bech32 Acc: ' | awk '{ print $NF }'`
arkeod tx arkeo open-contract --from $ark_user -- $ark_provider $ark_chain "$ark_pubkey" "$ark_contract_type" "$ark_deposit" "$ark_duration" $ark_rate "$ark_settle_duration"
```

## Step 4: Make Requests

There is a simple cli tool that was written. To install it, use `make tools`. The tool, `curleo`, abstracts away all of the signing and authentication elements for the user to simplify making an authenticated request.

```bash
 curleo -u bob -data '{ "jsonrpc": "2.0", "method": "health", "params": [], "id": 1 }' -H "text/plain" http://testnet-sentinel.arkeo.shapeshift.com:3636/gaia-mainnet-rpc-archive
invoking Sign...
Signed successfully
making POST request to http://testnet-sentinel.arkeo.shapeshift.com:3636/gaia-mainnet-rpc-archive?arkauth=7%3Atarkeopub1addwnpepqth6vxnuukr36du0cwz3cam63vqghgcs50wn9lxcsph3xddnq4y57538gg6%3A1%3Aadee652fad4de9de70a9d1246f3aedc4d0e9ef36197e48d3798961caa7a58bf346275d9cd57383a0abf8ad2621f9d6b908350399d62fe1b3dc1cad292063ab9c
{ "jsonrpc": "2.0", "method": "health", "params": [], "id": 1 }
{"jsonrpc":"2.0","id":1,"result":{}}
```

## Step 5: Claim Rewards for the Provider

On the provider’s behalf, you can claim the rewards for them (for testing purposes). To do use the following command. 

```bash
# nonce represents the number of requests made, must increase with each call for given contract
ark_nonce=20
ark_chain=gaia-mainnet-rpc-archive
ark_provider=$(curl -s http://directory.arkeo.shapeshift.com/provider/search/ | jq '.[]|select(.Status == "ONLINE" and .Chain == "gaia-mainnet-rpc-archive")|[{pubkey: .Pubkey, chain: .Chain, meta: .MetadataURI}]' | jq -r '.[0].pubkey')
ark_spender=$(arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r .key) | grep "Bech32 Acc" | awk '{ print $NF }')
ark_contract_id=$(arkeod query arkeo active-contract -o json $ark_spender $ark_provider $ark_chain | jq -r '.contract.id')
# if you get `rpc error: code = NotFound desc = not found: key not found` - it's likely there is no open contract
ark_sig=$(signhere -u $ark_user -m "$ark_contract_id:$ark_spender:$ark_nonce")
arkeod tx arkeo claim-contract-income --from $ark_user -- $ark_contract_id $ark_spender $ark_nonce $ark_sig
```

## Step 6: Close a Contract

If the contract is a subscription, it can be cancelled. Pay-as-you-go isn’t available to cancel as you can stop making requests as a form of cancelling (providers can cancel though). Closing a contract should also trigger a payout to the provider.

```bash
ark_chain=gaia-mainnet-rpc-archive
ark_spender=$(arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r .key) | grep "Bech32 Acc" | awk '{ print $NF }')
ark_contract_id=$(arkeod query arkeo active-contract -o json $ark_spender $ark_provider $ark_chain | jq -r '.contract.id')
arkeod tx arkeo close-contract --from $ark_user -- $contract_id
```
