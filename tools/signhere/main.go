package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// ModuleBasics is a mock module basic manager for testing
var ModuleBasics = module.NewBasicManager()

// main : Generate our pool address.
func main() {
	// network := flag.Int("n", 0, "The network to use.")
	user := flag.String("u", "alice", "user name")
	msg := flag.String("m", "message", "the text to sign")
	flag.Parse()

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	types.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	buf := bufio.NewReader(os.Stdin)

	kb, err := cKeys.New("mercuryd", cKeys.BackendTest, "~/.mercury", buf, cdc)
	if err != nil {
		log.Fatalf("%v", err)
	}

	bites := []byte(*msg)

	signature, pk, err := kb.Sign(*user, bites)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// verify signature
	if !pk.VerifySignature(bites, signature) {
		log.Fatal("bad signature")
	}

	fmt.Println(hex.EncodeToString(signature[:]))
}
