# Testnet Documentation Round

1. Prerequisites
    - a macOS or linux computer (docker can be used for windows)
    - a command shell and basic working knowledge
    - jq (used in example commands). install via package manager or [source](https://github.com/stedolan/jq).  
        macOS:
        ```bash
        brew install jq
        ```
        debian:
        ```bash
        apt-get install jq
        ```
1. Install client binaries  
    Either clone and build arkeo and its tools from source as outlined in [readme.md](../readme.md),
  or download a binary for your operating system below:
    - [macOS (x86-64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_darwin_amd64.tar.gz) 
      | [sha256](sums/arkeo_darwin_amd64.sha256?raw=true)
    - [macOS (arm64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_darwin_arm64.tar.gz) 
      | [sha256](sums/arkeo_darwin_arm64.sha256?raw=true)
    - [linux (x86-64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_linux_arm64.tar.gz)
      | [sha256](sums/arkeo_linux_arm64.sha256?raw=true)
    - [linux (arm64)](https://arkeo.s3.eu-west-1.amazonaws.com/bin/arkeo_linux_arm64.tar.gz)
      | [sha256](sums/arkeo_linux_arm64.sha256?raw=true)

    after downloading the executable and corresponding checksum file (right click->save as/link/target), verify the integrity of the downloaded artifact:  
      **Note**: `the examples are using macOS arm64 as the example, substitute your platform as needed.`
      
      the below will extract arkeod, curleo, and signhere
    to /usr/local/bin. replace "-C /usr/local/bin" with a different path if desired.  

    change your working directory to the directory you downloaded the archive and checksum files to.
    ```bash
    cd /path/to/downloads
    ```
    use the `sha256sum` command to check the archive against the checksum file
    ```bash
    sha256sum -c arkeo_darwin_arm64.sha256
    # Output:
    arkeo_darwin_arm64.tar.gz: OK
    ```
    note the output "`arkeo_darwin_arm64.tar.gz: OK`"
    
    extract the binaries from the archive.
    ```bash
    # add sudo if needed
    tar -C /usr/local/bin -zxvf arkeo_darwin_arm64.tar.gz
    # Output
    curleo
    arkeod
    signhere
    ```
    if the path you extracted the binaries to in the prior step is not on your PATH (/usr/local/bin generally is on *nix), add the directory to your PATH:
    ```bash
    export PATH=$PATH:/path/to/directory_containing
    ```
    /path/to/directory is the directory you extracted the binaries you downloaded to with the tar command above. to make this permanent,
    add the prior `export PATH=...` statement to your shell initialization scripts (.bashrc, .zshrc, etc.)

    verify installation by executing the arkeod command:
    ```bash
    arkeod version --long -o json | jq '.commit'
    # Output
    "793e955ba9d8d49609bb96fda6c85b0676419b58"
    ```
    update your client config as follows:
    ```bash
    arkeod config chain-id arkeo
    ```
    ```
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
    -and query the current block height:
    ```
    arkeod query block | jq -r '.block.header.height'
    205160
    ```
1. Setup a wallet  
  Create or Import a wallet using the cli. Replace `adam` below with a name of your choice. add `--recover` if you have a
mnemonic you'd like to use.

    assign a name to the $ark_user variable. this will become the key/wallet/user name.
    ```bash
    ark_user=adam
    ```
    create a new wallet (keypair)
    ```bash
    arkeod keys add $ark_user
    ```
    `OR` recover an existing wallet from bip39 mnemonic
    ```bash
    arkeod keys add $ark_user --recover
    ```
    This will output your arkeo address and pubkey along with the name you selected, as well as the mnemonic if you opted to generate one. Save this somewhere safe.

    ```bash
    - address: tarkeo14q4qnm4tmkm9xuhjwu0vw0f8xy7ztexek44nw3
      name: adam
      pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"ApcnOwEDoO6smze46IPUgC/5bC8DohEpLJ9ZZnrKky0w"}'
      type: local
    ```

    In order to interact with arkeo providers, you will need to have the `Acc` (account) pubkey encoded bech32 with the standard prefix. Execute the command below to obtain it.
    ```bash
    arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r '.key') | grep '^Bech32 Acc: ' | awk '{ print $NF }'
    # Output
    tarkeopub1addwnpepq2tjwwcpqwswatymx7uw3q75sqhljmp0qw3pz2fvnavkv7k2jvknqk25q7j
    ```

1. Get Test Funds
  Reach out to the arkeo development team with your address (starts with `tarkeo1`) to request funds.

1. Open a Contract

    Get a list of Online providers from the Directory Service:
    ```bash
    curl -s http://directory.arkeo.shapeshift.com/provider/search/ | jq '.[]|select(.Status == "ONLINE")|[{pubkey: .Pubkey, chain: .Chain, meta: .MetadataURI}]'

    # Output
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
    define vars:
    ```bash
    ark_contract_type=0
    ark_deposit=200
    ark_duration=20
    ark_settle_duration=10
    ark_rate=10
    ```
    obtain your (spender's) pubkey:
    ```bash
    ark_pubkey=`arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r '.key') | grep '^Bech32 Acc: ' | awk '{ print $NF }'`
    ```
    open the contract by broadcasting an arkeo open-contract transaction to the chain:
    ```bash
    arkeod tx arkeo open-contract --from $ark_user -- $ark_provider $ark_chain "$ark_pubkey" "$ark_contract_type" "$ark_deposit" "$ark_duration" $ark_rate "$ark_settle_duration"
    ```
    if things went well, the last line of output will be the txhash:
    ```bash
    txhash: 0B0AB5F982BFBB50E7518B114E02AF37D6A08E4E555FB47048461A632B1D0AC3
    ```
    check the status of your open-contract tx:
    ```bash
    arkeod query tx -o json --type=hash 0B0AB5F982BFBB50E7518B114E02AF37D6A08E4E555FB47048461A632B1D0AC3 | jq '.code'
    # Output
    0
    ```
    if the output is anything besides `0`, the tx didn't complete successfully. check the tx's raw_log:
    ```bash
    arkeod query tx -o json --type=hash 813A6B9A761F5A26E32EAECE5AFE73ED9D383D61C289AF515BC8D438B522883E | jq '.raw_log'
    # Output
    "failed to execute message; message index: 0: expires in 7 blocks: contract is already open"
    ```
    
    Open a Pay-As-You-Go contract. This example opens a pay-as-you-go contract for 20 blocks at a rate of 20 arkeo, depositing 400 to cover the duration.
      define vars:
    ```bash
    ark_contract_type=1
    ark_deposit=400
    ark_duration=20
    ark_settle_duration=10
    ark_rate=20
    ```
    obtain your (spender's) pubkey:
    ```bash
    ark_pubkey=`arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r '.key') | grep '^Bech32 Acc: ' | awk '{ print $NF }'`
    ```
    open the contract by broadcasting an arkeo open-contract transaction to the chain:
    ```bash
    arkeod tx arkeo open-contract --from $ark_user -- $ark_provider $ark_chain "$ark_pubkey" "$ark_contract_type" "$ark_deposit" "$ark_duration" $ark_rate "$ark_settle_duration"
    ```
1. Make Requests  
Use arkeo's `curleo` command to subchain rpc requests to the GAIA node:

    ```bash
    curleo -u $ark_user \                                                                                                            
    -data '{ "jsonrpc": "2.0", "method": "health", "params": [], "id": 1 }' \
    http://testnet-sentinel.arkeo.shapeshift.com:3636/gaia-mainnet-rpc-archive
    
    # Output
    invoking Sign...
    Signed successfully
    making POST request to http://testnet-sentinel.arkeo.shapeshift.com:3636/gaia-mainnet-rpc-archive?arkauth=17%3Atarkeopub1addwnpepq2tjwwcpqwswatymx7uw3q75sqhljmp0qw3pz2fvnavkv7k2jvknqk25q7j%3A1%3Ac2cafe3f9986cd11175baa723fd6491f021b00d9aedd2f3ade110525586a9cfc622e060129b7da0b98e05acf2458bad266e713048f415d901ea28fa522d2d37d
    { "jsonrpc": "2.0", "method": "health", "params": [], "id": 1 }
    {"jsonrpc":"2.0","id":1,"result":{}}
    ```

1. Claim Rewards  
On the provider’s behalf, you can claim the rewards for them (for testing purposes). To do use the following command. 

    ```bash
    # nonce represents the number of requests made, must increase with each call for given contract
    ark_nonce=20
    ark_chain=gaia-mainnet-rpc-archive
    ```
    
    ```bash
    ark_provider=$(curl -s http://directory.arkeo.shapeshift.com/provider/search/ | jq '.[]|select(.Status == "ONLINE" and .Chain == "gaia-mainnet-rpc-archive")|[{pubkey: .Pubkey, chain: .Chain, meta: .MetadataURI}]' | jq -r '.[0].pubkey')
    ```
    obtain your (spender's) pubkey:
    ```bash
    ark_spender=$(arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r .key) | grep "Bech32 Acc" | awk '{ print $NF }')
    ```
    obtain the active contract id
    ```bash
    ark_contract_id=$(arkeod query arkeo active-contract -o json $ark_spender $ark_provider $ark_chain | jq -r '.contract.id')
    ```
    use arkeo's `signhere` command to sign the authorization which must accompany paid tier requests
    ```bash
    ark_sig=$(signhere -u $ark_user -m "$ark_contract_id:$ark_spender:$ark_nonce")
    arkeod tx arkeo claim-contract-income --from $ark_user -- $ark_contract_id $ark_spender $ark_nonce $ark_sig
    ```

1. Close a Contract
Subscription contracts can be cancelled. Pay-as-you-go isn’t available to cancel as you can stop making requests as a form of cancelling.  
assign chain and obtain our (spender's) pubkey
```bash
ark_chain=gaia-mainnet-rpc-archive
ark_spender=$(arkeod debug pubkey-raw $(arkeod keys show $ark_user -p | jq -r .key) | grep "Bech32 Acc" | awk '{ print $NF }')
```
obtain the active contract id
```bash
ark_contract_id=$(arkeod query arkeo active-contract -o json $ark_spender $ark_provider $ark_chain | jq -r '.contract.id')
```
broadcast the arkeo close-contract transaction to the chain:
```bash
arkeod tx arkeo close-contract --from $ark_user -- $contract_id
```
