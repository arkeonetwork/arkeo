# Arkeo Testnet Setup

## Prerequisites

### Install Go

1. **Update and configure system parameters:**
    ```bash
    sudo apt-get update
    sudo apt-get upgrade -y
    sudo apt-get install -y build-essential curl wget jq make gcc chrony git
    sudo su -c "echo 'fs.file-max = 65536' >> /etc/sysctl.conf"
    sudo sysctl -p
    ```

2. **Install Go:**
    ```bash
    sudo rm -rf /usr/local/.go
    wget https://go.dev/dl/go1.21.13.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.13.linux-amd64.tar.gz
    sudo cp /usr/local/go /usr/local/.go -r
    sudo rm -rf /usr/local/go
    ```

3. **Update environment variables to include Go:**
    ```bash
    cat <<'EOF' >>$HOME/.profile
    export GOROOT=/usr/local/.go
    export GOPATH=$HOME/go
    export GO111MODULE=on
    export PATH=$PATH:/usr/local/.go/bin:$HOME/go/bin

    export ARKEO_PORT="<your-port>"
    EOF
    source $HOME/.profile
    ```

4. **Verify Go installation:**
    ```bash
    go version
    ```
    Expected output: `go version go1.21.13 linux/amd64`

---

## Install and Configure Arkeo Binary

### 1. **Clone and install Arkeo:**
   ```bash
   git clone https://github.com/arkeonetwork/arkeo
   cd arkeo
   git checkout master
   TAG=testnet make install
   ```

### 2. **Configure the binary:**
   ```bash
   arkeod keys add <key-name>
   arkeod config set client node tcp://localhost:${ARKEO_PORT}57
   arkeod config set client keyring-backend test
   arkeod config set client chain-id arkeo-testnet-3
   arkeod init <your-custom-moniker> --chain-id arkeo-testnet-3
   sudo ufw allow ${ARKEO_PORT}56/tcp
   ```

---

## Download Genesis File

```bash
curl -s http://seed.innovationtheory.com:26657/genesis | jq '.result.genesis' > $HOME/.arkeo/config/genesis.json
```

---

## Configure App Settings

### **Pruning, Gas Prices, and Prometheus: **
   ```bash
   sed -i -e "s/^pruning *=.*/pruning = \"custom\"/" $HOME/.arkeo/config/app.toml
   sed -i -e "s/^pruning-keep-recent *=.*/pruning-keep-recent = \"100\"/" $HOME/.arkeo/config/app.toml
   sed -i -e "s/^pruning-interval *=.*/pruning-interval = \"50\"/" $HOME/.arkeo/config/app.toml
   sed -i 's|minimum-gas-prices =.*|minimum-gas-prices = "0.001uarkeo"|g' $HOME/.arkeo/config/app.toml
   ```

---

## Configure Seeds and Peers

### **Set Seeds and Peers:**
   ```bash
   SEEDS="9dfa5f2d19c1174baf5e597965394269e654f9b7@seed31.innovationtheory.com:26656"
   PEERS="bb761c984bd990f3055f412917396754cd22af7a@validator31.innovationtheory.com:26656,81e36f94351d47803b8e1e0d0ad2d2e8e14ed36b@validator32.innovationtheory.com:26656"
   sed -i -e "/^\[p2p\]/,/^\[/{s/^[[:space:]]*seeds *=.*/seeds = \"$SEEDS\"/}" \
          -e "/^\[p2p\]/,/^\[/{s/^[[:space:]]*persistent_peers *=.*/persistent_peers = \"$PEERS\"/}" $HOME/.arkeo/config/config.toml
   ```

---

## Configure Node Service

### **Create the service file:**
   ```bash
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

### **Reset the node:**
   ```bash
   arkeod tendermint unsafe-reset-all --home $HOME/.arkeo --keep-addr-book
   ```

---

## Start and Enable the Service

### **Reload and start the service:**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable arkeo
   sudo systemctl restart arkeo && sudo journalctl -fu arkeo
   ```

---

> For any interface or dashboard calculations:
       **1 arkeo = 100,000,000 uarkeo**