package main

import (
	"bufio"
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

	"arkeo/app"
	"arkeo/common"
	"arkeo/common/cosmos"
	"arkeo/sentinel"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types"
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

	auth := curl.sign(*user, metadata.Configuration.ProviderPubKey.String(), chain, spender, claim.Height, claim.Nonce+1)
	values.Add(sentinel.QueryArkAuth, auth)

	u.RawQuery = values.Encode()

	resp, err := curl.client.Get(u.String())
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
	types.RegisterInterfaces(interfaceRegistry)
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
	types.RegisterInterfaces(interfaceRegistry)
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
