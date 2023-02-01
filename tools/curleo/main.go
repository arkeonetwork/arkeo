package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ArkeoNetwork/arkeo-protocol/app"
	"github.com/ArkeoNetwork/arkeo-protocol/common"
	"github.com/ArkeoNetwork/arkeo-protocol/common/cosmos"
	"github.com/ArkeoNetwork/arkeo-protocol/sentinel"
	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// ModuleBasics is a mock module basic manager for testing
var ModuleBasics = module.NewBasicManager()

type Curl struct {
	client  http.Client
	baseURL string
}

// main : Generate our pool address.
func main() {
	// network := flag.Int("n", 0, "The network to use.")
	user := flag.String("u", "alice", "user name")
	data := flag.String("data", "", "POST data")
	head := flag.String("H", "", "header")
	flag.Parse()

	c := cosmos.GetConfig()
	c.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")

	uriRaw := os.Args[len(os.Args)-1]
	u, err := url.Parse(uriRaw)
	if err != nil {
		log.Fatal(err)
	}
	values := u.Query()

	parts := strings.Split(u.Path, "/")
	chain := parts[1]

	curl := Curl{
		client:  http.Client{Timeout: time.Duration(5) * time.Second},
		baseURL: fmt.Sprintf("%s://%s", u.Scheme, u.Host),
	}
	metadata := curl.parseMetadata()
	spender := curl.getSpender(*user)
	claim := curl.getClaim(metadata.Configuration.ProviderPubKey.String(), chain, spender)
	height := claim.Height
	if height == 0 {
		contract := curl.getContract(metadata.Configuration.ProviderPubKey.String(), chain, spender)
		height = contract.Height
	}

	auth := curl.sign(*user, metadata.Configuration.ProviderPubKey.String(), chain, spender, height, claim.Nonce+1)
	values.Add(sentinel.QueryArkAuth, auth)

	u.RawQuery = values.Encode()

	var resp *http.Response

	if len(*data) > 0 {
		header := "application/x-www-form-urlencoded"
		if len(*head) > 0 {
			header = *head
		}
		resp, err = curl.client.Post(u.String(), header, bytes.NewBuffer([]byte(*data)))
	} else {
		resp, err = curl.client.Get(u.String())
	}
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err) // nolint
	}

	fmt.Println(string(body))
}

func (c Curl) getContract(provider, chain, spender string) types.Contract {
	url := fmt.Sprintf("%s/contract/%s/%s/%s", c.baseURL, provider, chain, spender)
	resp, err := c.client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err) // nolint
	}

	var claim types.Contract
	err = json.Unmarshal(body, &claim)
	if err != nil {
		log.Fatal(err)
	}

	return claim
}

func (c Curl) getClaim(provider, chain, spender string) sentinel.Claim {
	url := fmt.Sprintf("%s/claim/%s/%s/%s", c.baseURL, provider, chain, spender)
	resp, err := c.client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err) // nolint
	}

	var claim sentinel.Claim
	err = json.Unmarshal(body, &claim)
	if err != nil {
		log.Fatal(err)
	}

	return claim
}

func (c Curl) parseMetadata() sentinel.Metadata {
	metadataURI := fmt.Sprintf("%s/metadata.json", c.baseURL)
	resp, err := c.client.Get(metadataURI)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err) // nolint
	}

	var meta sentinel.Metadata
	err = json.Unmarshal(body, &meta)
	if err != nil {
		log.Fatal(err) // nolint
	}

	return meta
}

func (c Curl) sign(user, provider, chain, spender string, height, nonce int64) string {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	sdk.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	buf := bufio.NewReader(os.Stdin)

	kb, err := cKeys.New("arkeod", cKeys.BackendTest, "~/.arkeo", buf, cdc)
	if err != nil {
		log.Fatal(err)
	}

	msg := fmt.Sprintf("%s:%s:%s:%d:%d", provider, chain, spender, height, nonce)

	signature, pk, err := kb.Sign(user, []byte(msg))
	if err != nil {
		log.Fatal(err)
	}

	// verify signature
	if !pk.VerifySignature([]byte(msg), signature) {
		log.Fatal("bad signature")
	}

	sig := hex.EncodeToString(signature)
	return fmt.Sprintf("%s:%s", msg, sig)
}

func (c Curl) getSpender(user string) string {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	sdk.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	buf := bufio.NewReader(os.Stdin)

	kb, err := cKeys.New("arkeod", cKeys.BackendTest, "~/.arkeo", buf, cdc)
	if err != nil {
		log.Fatal(err)
	}

	record, err := kb.Key(user)
	if err != nil {
		log.Fatal(err) // nolint
	}

	pub, err := record.GetPubKey()
	if err != nil {
		log.Fatal(err) // nolint
	}

	pk, err := common.NewPubKeyFromCrypto(pub)
	if err != nil {
		log.Fatal(err) // nolint
	}

	return pk.String()
}
