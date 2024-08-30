# Testnet Setup

 Arkeo Testnet Setup

## Prerequisites 

### Install Go 

Make sure your system is updated and set the system parameters correctly

```shell
sudo apt-get update
sudo apt-get upgrade -y
sudo apt-get install -y build-essential curl wget jq make gcc chrony git
sudo su -c "echo 'fs.file-max = 65536' >> /etc/sysctl.conf"
sudo sysctl -p
```
Install GO 

```shell
sudo rm -rf /usr/local/.go
wget https://go.dev/dl/go1.21.13.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.13.linux-amd64.tar.gz
sudo cp /usr/local/go /usr/local/.go -r 
sudo rm -rf /usr/local/go
```
Update environment variables to include go

```shell
cat <<'EOF' >>$HOME/.profile
export GOROOT=/usr/local/.go
export GOPATH=$HOME/go
export GO111MODULE=on
export PATH=$PATH:/usr/local/.go/bin:$HOME/go/bin

export ARKEO_PORT="<your-port>"
EOF
source $HOME/.profile
```

Check if go is correctly installed
```shell
go version 
```

This should return something like `go version go1.21.13 linux/amd64`

# Arkeo Binary 

Install the Arkeo Binary 
```shell
git clone https://github.com/arkeonetwork/arkeo
cd arkeo
git checkout master
TAG=testnet make install 
```

Configure The Binary 

```shell

arkeod keys add <key-name>
arkeod config set client node tcp://localhost:${ARKEO_PORT}57
arkeod config set client keyring-backend os
arkeod config set client chain-id arkeo
arkeod init <your-custom-moniker> --chain-id arkeo
sudo ufw allow ${ARKEO_PORT}56/tcp
```

## Download Genesis 
```shell
curl -s http://seed.innovationtheory.com:26657/genesis | jq '.result.genesis' > $HOME/.arkeo/config/genesis.json
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
SEEDS="aab68f68841eb072d996cd1b45c2b9c9b612d95b@seed.innovationtheory.com:26656,85341b428cf5993fcc04a324d95d14590ae5172c@eed2.innovationtheory.com:26656"
PEERS="c27c96c5b54a9f2bea776858e2cff364e410d2a8@71.218.54.128:26656,f7da702c17e45e463adf21e57b1d0d936cbc97a3@peer2.innovationtheory.com:26656,46e6d4751bbc67d3e72e13dacdfb0770227fbfc3@65.108.79.241:46656,fc5464b2ce731c5787be0fd316b6c4b6611886ea@37.252.184.241:26656,
57f693ba3fed4dd82d02d4cbcc73712c6da4bd34@65.109.113.228:60756,38ab548031dea2b46889c253b762e51306d29cbb@65.109.92.148:61656"
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

Reset Node
```shell
arkeod tendermint unsafe-reset-all --home $HOME/.arkeo --keep-addr-book
```

Enable and start the node service

```shell
sudo systemctl daemon-reload
sudo systemctl enable arkeo
sudo systemctl restart arkeo && sudo journalctl -fu arkeo
```