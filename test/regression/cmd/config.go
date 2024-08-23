package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	tmhttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	arekoappParams "github.com/arkeonetwork/arkeo/app/params"
)

////////////////////////////////////////////////////////////////////////////////////////
// Cosmos
////////////////////////////////////////////////////////////////////////////////////////

// TODO: determine how to return these programmatically without keeper
const (
	ModuleAddr         = "tarkeo1k8925g52vwe5jgfp4nqr7ljuhs7nzu2f8g0za5"
	ModuleAddrContract = "tarkeo1kz2dkl8zlxwte008astc5e65htrxdcse6x3h3h"
	ModuleAddrProvider = "tarkeo1rcvm4v5mcepj53fh2526uve0tly4grdsx5yw7k"
	ModuleAddrReserve  = "tarkeo1d0m97ywk2y4vq58ud6q5e0r3q9khj9e3unfe4t"
)

var (
	encodingConfig arekoappParams.EncodingConfig
	clientCtx      client.Context
	txFactory      tx.Factory
	keyRing        keyring.Keyring
)

func init() {
	// initialize the bech32 prefix for testnet/mocknet
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")
	config.Seal()

	// initialize the codec
	encodingConfig = app.MakeEncodingConfig()

	// create new rpc client
	rpcClient, err := tmhttp.New("http://localhost:26657", "/websocket")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create tendermint client")
	}

	// create keyring
	keyRing = keyring.NewInMemory(encodingConfig.Marshaler)

	// create cosmos-sdk client context
	clientCtx = client.Context{
		Client:            rpcClient,
		ChainID:           "arkeo",
		Codec:             encodingConfig.Marshaler,
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		Keyring:           keyRing,
		BroadcastMode:     flags.BroadcastSync,
		SkipConfirm:       true,
		TxConfig:          encodingConfig.TxConfig,
		AccountRetriever:  authtypes.AccountRetriever{},
		NodeURI:           "http://localhost:26657",
		LegacyAmino:       encodingConfig.Amino,
	}

	// create tx factory
	txFactory = txFactory.WithKeybase(clientCtx.Keyring)
	txFactory = txFactory.WithTxConfig(clientCtx.TxConfig)
	txFactory = txFactory.WithAccountRetriever(clientCtx.AccountRetriever)
	txFactory = txFactory.WithChainID(clientCtx.ChainID)
	txFactory = txFactory.WithGas(1e8)
	txFactory = txFactory.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
	txFactory = txFactory.WithTxConfig(clientCtx.TxConfig)
}

////////////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////////////

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	// set to info level if DEBUG is not set (debug is the default level)
	if os.Getenv("DEBUG") == "" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Colors
////////////////////////////////////////////////////////////////////////////////////////

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorPurple = "\033[35m"

	// save for later
	// ColorYellow = "\033[33m"
	// ColorBlue   = "\033[34m"
	// ColorCyan   = "\033[36m"
	// ColorGray   = "\033[37m"
	// ColorWhite  = "\033[97m"
)

////////////////////////////////////////////////////////////////////////////////////////
// HTTP
////////////////////////////////////////////////////////////////////////////////////////

var httpClient = &http.Client{
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second * getTimeFactor(),
		}).Dial,
	},
	Timeout: 30 * time.Second * getTimeFactor(),
}

// //////////////////////////////////////////////////////////////////////////////////////
// Module Addresses
// //////////////////////////////////////////////////////////////////////////////////////
// nolint
// trunk-ignore-all(gitleaks/generic-api-key)
// trunk-ignore(trunk/ignore-does-nothing) # Getting weird lint error only on CI and running locally works fine
// trunk-ignore-all(golangci-lint/gosec)
const (
	ModuleAddrBondedTokensPool    = "tarkeo1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3e79s43" //nolint:staticcheck
	ModuleAddrNotBondedTokensPool = "tarkeo1tygms3xhhs3yv487phx3dw4a95jn7t7ld7epr9" //nolint:staticcheck
	ModuleAddrGov                 = "tarkeo10d07y265gmmuvt4z0w9aw880jnsr700jk8l664" //nolint:staticcheck
	ModuleAddrDistribution        = "tarkeo1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8t6gr9e" //nolint:staticcheck
	ModuleAddrMint                = "tarkeo1m3h30wlvsf8llruxtpukdvsy0km2kum8y5t8tx" //nolint:staticcheck
	ModuleAddrFeeCollector        = "tarkeo17xpfvakm2amg962yls6f84z3kell8c5luu0l8m" //nolint:staticcheck
)

////////////////////////////////////////////////////////////////////////////////////////
// Keys
////////////////////////////////////////////////////////////////////////////////////////

var (
	addressToName   = map[string]string{} // arkeo...->dog, 0x...->dog
	templateAddress = map[string]string{} // addr_dog->arkeo...
	templatePubKey  = map[string]string{} // pubkey_dog->arkeopub...

	dogMnemonic = strings.Repeat("dog ", 23) + "fossil"
	catMnemonic = strings.Repeat("cat ", 23) + "crawl"
	foxMnemonic = strings.Repeat("fox ", 23) + "filter"
	pigMnemonic = strings.Repeat("pig ", 23) + "quick"

	// mnemonics contains the set of all mnemonics for accounts used in tests
	mnemonics = [...]string{
		dogMnemonic,
		catMnemonic,
		foxMnemonic,
		pigMnemonic,
	}
)

func init() {
	// get default hd path
	hdPath := hd.NewFundraiserParams(0, 118, 0).String()

	// register functions for all mnemonic-chain addresses
	for _, m := range mnemonics {
		name := strings.Split(m, " ")[0]

		// create pubkey for mnemonic
		derivedPriv, err := hd.Secp256k1.Derive()(m, "", hdPath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to derive private key")
		}
		privKey := hd.Secp256k1.Generate()(derivedPriv)
		s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, privKey.PubKey())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to bech32ify pubkey")
		}
		pk := common.PubKey(s)

		// add key to keyring
		_, err = keyRing.NewAccount(name, m, "", hdPath, hd.Secp256k1)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to add account to keyring")
		}

		// register template address
		addr, err := pk.GetMyAddress()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get address")
		}
		templateAddress[fmt.Sprintf("addr_%s", name)] = addr.String()

		// register address to name
		addressToName[addr.String()] = name

		// register pubkey for arkeonetwork
		templatePubKey[fmt.Sprintf("pubkey_%s", name)] = pk.String()
	}
}
