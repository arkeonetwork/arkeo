# Testnet Setup

Run testnet setup using the pre-complied testnet binary.

## Download Testnet Binary
```shell
./scripts/install-binary.sh
```

## Build Your Own Binary 

This installs binary to /usr/local/bin
```shell
make install-testnet-binary 
```


## Configure variables 

Set variables

```shell
echo "export ARKEO_WALLET="MyWallet"" >> $HOME/.bash_profile
echo "export MONIKER="MyNode"" >> $HOME/.bash_profile
echo "export ARKEO_CHAIN_ID="arkeo"" >> $HOME/.bash_profile
echo "export ARKEO_PORT="10"" >> $HOME/.bash_profile


source $HOME/.bash_profile
```



Configure `client.toml` 
```shell
arkeod config set client node tcp://localhost:${ARKEO_PORT}57
arkeod config set client keyring-backend os
arkeod config set client chain-id arkeo 
```

Init App 
```shell
arkeod init $MONIKER --chain-id arkeo
```

## Download Genesis 
```shell
wget -qO- http://seed.innovationtheory.com:26657/genesis | jq '.result.genesis' > $HOME/.arkeo/config/genesis.json
```
## Set Custom Ports

In `app.toml`
```shell
sed -i.bak -e "s%:1317%:${ARKEO_PORT}17%g;
s%:8080%:${ARKEO_PORT}80%g;
s%:9090%:${ARKEO_PORT}90%g;
s%:9091%:${ARKEO_PORT}91%g;
s%:8545%:${ARKEO_PORT}45%g;
s%:8546%:${ARKEO_PORT}46%g;
s%:6065%:${ARKEO_PORT}65%g" $HOME/.arkeo/config/app.toml
```
In `config.toml`

```shell
sed -i.bak -e "s%:26658%:${ARKEO_PORT}58%g;
s%:26657%:${ARKEO_PORT}57%g;
s%:6060%:${ARKEO_PORT}60%g;
s%:26656%:${ARKEO_PORT}56%g;
s%^external_address = \"\"%external_address = \"$(wget -qO- eth0.me):${ARKEO_PORT}56\"%;
s%:26660%:${ARKEO_PORT}61%g" $HOME/.arkeo/config/config.toml
```

## Configure Pruning, Minimum gas price , enable prometheus and disable indexing 

```shell

sed -i -e "s/^pruning *=.*/pruning = \"custom\"/" $HOME/.arkeo/config/app.toml
sed -i -e "s/^pruning-keep-recent *=.*/pruning-keep-recent = \"100\"/" $HOME/.arkeo/config/app.toml
sed -i -e "s/^pruning-interval *=.*/pruning-interval = \"50\"/" $HOME/.arkeo/config/app.toml

sed -i 's|minimum-gas-prices =.*|minimum-gas-prices = "0.001uarkeo"|g' $HOME/.arkeo/config/app.toml
sed -i -e "s/prometheus = false/prometheus = true/" $HOME/.arkeo/config/config.toml
sed -i -e "s/^indexer *=.*/indexer = \"null\"/" $HOME/.arkeo/config/config.toml

```

## Configure Seeds and Peers

```shell
SEEDS="aab68f68841eb072d996cd1b45c2b9c9b612d95b@seed.innovationtheory.com:26656"
PEERS="dc824a97b930d9c8086d96ea4da515ffccc050d6@23.22.39.126:26656"
sed -i -e "/^\[p2p\]/,/^\[/{s/^[[:space:]]*seeds *=.*/seeds = \"$SEEDS\"/}" \
       -e "/^\[p2p\]/,/^\[/{s/^[[:space:]]*persistent_peers *=.*/persistent_peers = \"$PEERS\"/}" $HOME/.arkeo/config/config.toml
```

# Configure Service

Create the node service file

```shell
sudo tee /etc/systemd/system/arkeo.service > /dev/null <<EOF
[Unit]
Description=Arkeo node
After=network-online.target
[Service]
User=$USER
WorkingDirectory=$HOME/.arkeo
ExecStart= $(which arkeod) start --home $HOME/.arkeo
Restart=on-failure
RestartSec=5
LimitNOFILE=65535
[Install]
WantedBy=multi-user.target
EOF
```

Reset node and download snapshot
```shell
arkeod tendermint unsafe-reset-all --home $HOME/.arkeo --keep-addr-book
```

Allow node P2P port in firewall

```shell
sudo ufw allow ${ARKEO_PORT}56/tcp
```

Enable and start the node service

```shell
sudo systemctl daemon-reload
sudo systemctl enable arkeo
sudo systemctl restart arkeo && sudo journalctl -fu arkeo
```
