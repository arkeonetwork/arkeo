# Setting Up Testnet using Cosmovisor 

## How does cosmovisor work?
Cosmovisor is designed to be used as an abstract interface for the Cosmos SDK chain, for example Arkeo. Some key takeaways about Cosmovisor:
- It passes arguments to arkeod (which is configured by DAEMON_NAME env variable).
- It manages arkeod by restarting and upgrading if needed.
- It is configured using environment variables, not positional arguments.
- Running cosmovisor run arg1 arg2 .... runs arkeod arg1 arg2 ...
- All arguments passed to cosmovisor run are passed to the application binary, as a subprocess. As cosmovisor returns /dev/stdout and /dev/stderr of the subprocess as its own, cosmovisor run cannot accept any command-line arguments other than those available to arkeod.

Make sure to check the [cosmovisor documentation](https://docs.cosmos.network/main/tooling/cosmovisor) for a comprehensive guide on how to use Cosmovisor.

## Install Cosmovisor 

- Install [Prerequisites](./TESTNET.md#prerequisites) , Configure [Binary](./TESTNET.md#arkeo-binary) before installing `Cosmovisor`

Install latest Cosmovisor 

```shell
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
```

Verify the installation 

```shell
cosmovisor version
```

## Configure Cosmovisor and Environment Variables 
```shell
mkdir -p "${HOME}"/.arkeo/cosmovisor/genesis/bin
mkdir "${HOME}"/.arkeo/cosmovisor/upgrades
```

Then copy the arkeod binary to `genesis/bin` directory 
```shell
cp "${GOPATH}"/bin/arkeod "${HOME}"/.arkeo/cosmovisor/genesis/bin
```
Add the necessary environment variables, for example by adding these variables to the profile that will be running Cosmovisor. You can edit the ~/.profile file by adding the following content:

```shell
export DAEMON_HOME="${HOME}"/.arkeo
export DAEMON_RESTART_AFTER_UPGRADE=true
export DAEMON_ALLOW_DOWNLOAD_BINARIES=false
export DAEMON_NAME=arkeod
export UNSAFE_SKIP_BACKUP=true
```
Before running cosmovisor , you should make sure to have initialized the arkeo node and configured joining network by downloading genesis file , setting up peers and snapshot if any.

You can create a service file with:
```shell
sudo nano /etc/systemd/system/cosmovisor.service
```

and add the following content by making sure to change the <your-user>, <path-to-cosmovisor> and <path-to-arkeo> with your values:

```shell
[Unit]
Description=cosmovisor
After=network-online.target
[Service]
User=<your-user>
ExecStart=cosmovisor run  start --x-crisis-skip-assert-invariants
Restart=always
RestartSec=3
LimitNOFILE=4096
Environment="DAEMON_NAME=arkeod"
Environment="DAEMON_HOME=/<path-to-arkeo>/.arkeo"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="DAEMON_LOG_BUFFER_SIZE=512"
Environment="UNSAFE_SKIP_BACKUP=true"
[Install]
WantedBy=multi-user.target
```

You can now reload the systemctl daemon:
```shell
sudo -S systemctl daemon-reload
```

and enable Cosmovisor as a service:
```shell
sudo -S systemctl enable cosmovisor
```

You can now start Cosmovisor by executing:
```shell
sudo systemctl start cosmovisor
```

Make sure to check that the service is running by executing:
```shell
sudo systemctl status cosmovisor
```