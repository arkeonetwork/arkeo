package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/arkeonetwork/arkeo/sentinel"
	arkeo "github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

////////////////////////////////////////////////////////////////////////////////////////
// Operation
////////////////////////////////////////////////////////////////////////////////////////

type Operation interface {
	Execute(arkeo *os.Process, logs chan string) error
	OpType() string
}

type OpBase struct {
	Type string `json:"type"`
}

func (op *OpBase) OpType() string {
	return op.Type
}

func NewOperation(opMap map[string]any) Operation {
	// ensure type is provided
	t, ok := opMap["type"].(string)
	if !ok {
		log.Fatal().Interface("type", opMap["type"]).Msg("operation type is not a string")
	}

	// create the operation for the type
	var op Operation
	switch t {
	case "state":
		op = &OpState{}
	case "check":
		op = &OpCheck{}
	case "create-blocks":
		op = &OpCreateBlocks{}
	case "tx-send":
		op = &OpTxSend{}
	case "tx-bond-provider":
		op = &OpTxBondProvider{}
	case "tx-mod-provider":
		op = &OpTxModProvider{}
	case "tx-open-contract":
		op = &OpTxOpenContract{}
	case "tx-close-contract":
		op = &OpTxCloseContract{}
	case "tx-claim-contract":
		op = &OpTxClaimContract{}
	default:
		log.Fatal().Str("type", t).Msg("unknown operation type")
	}

	// create decoder supporting embedded structs and weakly typed input
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		ErrorUnused:      true,
		Squash:           true,
		Result:           op,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create decoder")
	}

	switch op.(type) {
	// internal types have MarshalJSON methods necessary to decode
	case *OpCheck, *OpTxSend, *OpTxBondProvider, *OpTxModProvider, *OpTxOpenContract, *OpTxCloseContract, *OpTxClaimContract:
		// encode as json
		buf := bytes.NewBuffer(nil)
		enc := json.NewEncoder(buf)
		err = enc.Encode(opMap)
		if err != nil {
			log.Fatal().Interface("op", opMap).Err(err).Msg("failed to encode operation")
		}

		// unmarshal json to op
		err = json.NewDecoder(buf).Decode(op)

	default:
		err = dec.Decode(opMap)
	}
	if err != nil {
		log.Fatal().Interface("op", opMap).Err(err).Msg("failed to decode operation")
	}

	// require check description and default status check to 200 if endpoint is set
	if oc, ok := op.(*OpCheck); ok && oc.Endpoint != "" {
		if oc.Description == "" {
			log.Fatal().Interface("op", opMap).Msg("check operation must have a description")
		}
		if oc.Status == 0 {
			oc.Status = 200
		}
	}

	return op
}

////////////////////////////////////////////////////////////////////////////////////////
// OpState
////////////////////////////////////////////////////////////////////////////////////////

type OpState struct {
	OpBase  `yaml:",inline"`
	Genesis map[string]any `json:"genesis"`
}

func (op *OpState) Execute(*os.Process, chan string) error {
	// load genesis file
	f, err := os.OpenFile(os.ExpandEnv("/regtest/.arkeo/config/genesis.json"), os.O_RDWR, 0o644)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open genesis file")
	}

	// unmarshal genesis into map
	var genesisMap map[string]any
	err = json.NewDecoder(f).Decode(&genesisMap)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to decode genesis file")
	}

	// merge updates into genesis
	genesis := deepMerge(genesisMap, op.Genesis)

	// reset file
	err = f.Truncate(0)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to truncate genesis file")
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to seek genesis file")
	}

	// marshal genesis into file
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(genesis)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to encode genesis file")
	}

	return f.Close()
}

////////////////////////////////////////////////////////////////////////////////////////
// OpCheck
////////////////////////////////////////////////////////////////////////////////////////

type OpCheck struct {
	OpBase        `yaml:",inline"`
	Description   string            `json:"description"`
	Endpoint      string            `json:"endpoint"`
	Method        string            `json:"method"`
	Body          string            `json:"body"`
	Params        map[string]string `json:"params"`
	ArkAuth       map[string]string `json:"arkauth"`
	ContractAuth  map[string]string `json:"contractauth"`
	Status        int               `json:"status"`
	AssertHeaders map[string]string `json:"headers"`
	Asserts       []string          `json:"asserts"`
}

func createContractAuth(input map[string]string) (string, bool, error) {
	///////// create contract auth signature //////////////
	if len(input) == 0 {
		return "", false, nil
	}

	// validate inputs
	if len(input["id"]) == 0 {
		return "", true, fmt.Errorf("missing required field: id")
	}

	if len(input["timestamp"]) == 0 {
		return "", true, fmt.Errorf("missing required field: timestamp")
	}

	id, err := strconv.ParseUint(input["id"], 10, 64)
	if err != nil {
		return "", true, fmt.Errorf("failed to parse id: %s", err)
	}
	timestamp, err := strconv.ParseInt(input["timestamp"], 10, 64)
	if err != nil {
		return "", true, fmt.Errorf("failed to parse timestamp: %s", err)
	}

	// sign our msg
	msg := sentinel.GenerateMessageToSign(id, timestamp)
	auth, err := signThis(msg, input["signer"])
	return auth, true, err
}

func createAuth(input map[string]string) (string, bool, error) {
	///////// create ark auth signature //////////////
	if len(input) == 0 {
		return "", false, nil
	}

	// validate inputs
	if len(input["id"]) == 0 {
		return "", true, fmt.Errorf("missing required field: id")
	}

	if _, ok := input["nosig"]; ok {
		// no signature requested
		return input["id"], true, nil
	}

	if len(input["nonce"]) == 0 {
		return "", true, fmt.Errorf("missing required field: nonce")
	}

	id, err := strconv.ParseUint(input["id"], 10, 64)
	if err != nil {
		return "", true, fmt.Errorf("failed to parse id: %s", err)
	}
	nonce, err := strconv.ParseInt(input["nonce"], 10, 64)
	if err != nil {
		return "", true, fmt.Errorf("failed to parse nonce: %s", err)
	}
	msg := sentinel.GenerateMessageToSign(id, nonce)
	auth, err := signThis(msg, input["signer"])
	return auth, true, err
}

func signThis(msg, signer string) (string, error) {
	///////// create ark auth signature //////////////
	mnemonic := ""
	switch signer {
	case "dog":
		mnemonic = dogMnemonic
	case "cat":
		mnemonic = catMnemonic
	case "fox":
		mnemonic = foxMnemonic
	case "pig":
		mnemonic = pigMnemonic
	default:
		return "", fmt.Errorf("ark auth requires a valid signer (dog, cat, fox, pig): %s", signer)
	}
	// get default hd path
	hdPath := hd.NewFundraiserParams(0, 118, 0).String()

	// create pubkey for mnemonic
	derivedPriv, err := hd.Secp256k1.Derive()(mnemonic, "", hdPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to derive private key")
	}
	privKey := hd.Secp256k1.Generate()(derivedPriv)

	sig, err := privKey.Sign([]byte(msg))
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %s", err)
	}
	return fmt.Sprintf("%s:%s", msg, hex.EncodeToString(sig)), nil
}

func (op *OpCheck) Execute(_ *os.Process, logs chan string) error {
	// abort if no endpoint is set (empty check op is allowed for breakpoint convenience)
	if op.Endpoint == "" {
		return fmt.Errorf("check")
	}

	if op.Method == "" {
		op.Method = http.MethodGet
	}

	var body io.Reader
	if len(op.Body) > 0 {
		body = strings.NewReader(op.Body)
	}

	// build request
	req, err := http.NewRequest(op.Method, op.Endpoint, body)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build request")
	}

	// add params
	q := req.URL.Query()
	for k, v := range op.Params {
		q.Add(k, v)
	}
	arkauth, authOK, err := createAuth(op.ArkAuth)
	if err != nil {
		return err
	}
	if authOK {
		q.Add(sentinel.QueryArkAuth, arkauth)
	}
	contractAuth, contractAuthOK, err := createContractAuth(op.ContractAuth)
	if err != nil {
		return err
	}
	if contractAuthOK {
		q.Add(sentinel.QueryContract, contractAuth)
	}
	req.URL.RawQuery = q.Encode()

	// send request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Err(err).Msg("failed to send request")
		return err
	}

	// read response
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Msg("failed to read response")
		return err
	}

	// ensure status code matches
	if op.Status == 0 { // default to 200
		op.Status = 200
	}
	if resp.StatusCode != op.Status {
		// dump pretty output for debugging
		fmt.Println(ColorPurple + "\nOperation:" + ColorReset)
		_ = yaml.NewEncoder(os.Stdout).Encode(op)
		fmt.Println(ColorPurple + "\nEndpoint Response:" + ColorReset)
		fmt.Println(string(buf) + "\n")

		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	for k, v := range op.AssertHeaders {
		if val, exists := resp.Header[k]; exists {
			if !strings.EqualFold(val[0], v) {
				return fmt.Errorf("Bad header: expected %s, got %s", v, val[0])
			}
		} else {
			return fmt.Errorf("Missing header: %s", k)
		}
	}

	// pipe response to jq for assertions
	for _, a := range op.Asserts {
		// render the assert expression (used for native_txid)
		tmpl := template.Must(template.Must(templates.Clone()).Parse(a))
		expr := bytes.NewBuffer(nil)
		err = tmpl.Execute(expr, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to render assert expression")
		}
		a = expr.String()

		cmd := exec.Command("jq", "-e", a)
		cmd.Stdin = bytes.NewReader(buf)
		out, err := cmd.CombinedOutput()
		if err != nil {
			if cmd.ProcessState.ExitCode() == 1 {
				// dump process logs if the assert expression failed
				fmt.Println(ColorPurple + "\nLogs:" + ColorReset)
				dumpLogs(logs)
			}

			// dump pretty output for debugging
			fmt.Println(ColorPurple + "\nOperation:" + ColorReset)
			_ = yaml.NewEncoder(os.Stdout).Encode(op)
			fmt.Println(ColorPurple + "\nFailed Assert: " + ColorReset + expr.String())
			fmt.Println(ColorPurple + "\nEndpoint Response:" + ColorReset)
			fmt.Println(string(buf) + "\n")

			// log fatal on syntax errors and skip logs
			if cmd.ProcessState.ExitCode() != 1 {
				drainLogs(logs)
				fmt.Println(ColorRed + string(out) + ColorReset)
			}

			return err
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////
// OpCreateBlocks
////////////////////////////////////////////////////////////////////////////////////////

type OpCreateBlocks struct {
	OpBase `yaml:",inline"`
	Count  int  `json:"count"`
	Exit   *int `json:"exit"`
}

func (op *OpCreateBlocks) Execute(p *os.Process, logs chan string) error {
	// clear existing log output
	drainLogs(logs)

	for i := 0; i < op.Count; i++ {
		// http request to localhost to unblock block creation
		_, err := httpClient.Get("http://localhost:8080/newBlock")
		if err != nil {
			// if exit code is not set this was unexpected
			if op.Exit == nil {
				log.Err(err).Msg("failed to create block")
				return err
			}

			// if exit code is set, this was expected
			if processRunning(p.Pid) {
				log.Err(err).Msg("block did not exit as expected")
				return err
			}

			// if process is not running, check exit code
			ps, err := p.Wait()
			if err != nil {
				log.Err(err).Msg("failed to wait for process")
				return err
			}
			if ps.ExitCode() != *op.Exit {
				log.Error().Int("exit", ps.ExitCode()).Int("expect", *op.Exit).Msg("bad exit code")
				return err
			}

			// exit code is correct, return nil
			return nil
		}
	}

	// if exit code is set, this was unexpected
	if op.Exit != nil {
		log.Error().Int("expect", *op.Exit).Msg("expected exit code")
		return errors.New("expected exit code")
	}

	// avoid minor raciness after end block
	time.Sleep(200 * time.Millisecond * getTimeFactor())

	return nil
}

// ------------------------------ OpTxSend ------------------------------

type OpTxSend struct {
	OpBase       `yaml:",inline"`
	bank.MsgSend `yaml:",inline"`
	Sequence     *int64 `json:"sequence"`
}

func (op *OpTxSend) Execute(_ *os.Process, logs chan string) error {
	signer := sdk.MustAccAddressFromBech32(op.FromAddress)
	return sendMsg(&op.MsgSend, signer, op.Sequence, op, logs)
}

// ------------------------------ OpTxBondProvider ------------------------------

type OpTxBondProvider struct {
	OpBase                `yaml:",inline"`
	arkeo.MsgBondProvider `yaml:",inline"`
	Signer                string `json:"signer"`
	Sequence              *int64 `json:"sequence"`
}

func (op *OpTxBondProvider) Execute(_ *os.Process, logs chan string) error {
	signer := sdk.MustAccAddressFromBech32(op.Signer)
	return sendMsg(&op.MsgBondProvider, signer, op.Sequence, op, logs)
}

// ------------------------------ OpTxModProvider ------------------------------

type OpTxModProvider struct {
	OpBase               `yaml:",inline"`
	arkeo.MsgModProvider `yaml:",inline"`
	Signer               string `json:"signer"`
	Sequence             *int64 `json:"sequence"`
}

func (op *OpTxModProvider) Execute(_ *os.Process, logs chan string) error {
	signer := sdk.MustAccAddressFromBech32(op.Signer)
	return sendMsg(&op.MsgModProvider, signer, op.Sequence, op, logs)
}

// ------------------------------ OpTxOpenContract ------------------------------

type OpTxOpenContract struct {
	OpBase                `yaml:",inline"`
	arkeo.MsgOpenContract `yaml:",inline"`
	Signer                string `json:"signer"`
	Sequence              *int64 `json:"sequence"`
}

func (op *OpTxOpenContract) Execute(_ *os.Process, logs chan string) error {
	signer := sdk.MustAccAddressFromBech32(op.Signer)
	return sendMsg(&op.MsgOpenContract, signer, op.Sequence, op, logs)
}

// ------------------------------ OpTxCloseContract ------------------------------

type OpTxCloseContract struct {
	OpBase                 `yaml:",inline"`
	arkeo.MsgCloseContract `yaml:",inline"`
	Signer                 string `json:"signer"`
	Sequence               *int64 `json:"sequence"`
}

func (op *OpTxCloseContract) Execute(_ *os.Process, logs chan string) error {
	signer := sdk.MustAccAddressFromBech32(op.Signer)
	return sendMsg(&op.MsgCloseContract, signer, op.Sequence, op, logs)
}

// ------------------------------ OpTxClaimContract ------------------------------

type OpTxClaimContract struct {
	OpBase                       `yaml:",inline"`
	arkeo.MsgClaimContractIncome `yaml:",inline"`
	Signer                       string            `json:"signer"`
	ArkAuth                      map[string]string `json:"arkauth"`
	Sequence                     *int64            `json:"sequence"`
}

func (op *OpTxClaimContract) Execute(_ *os.Process, logs chan string) error {
	var err error
	arkauth, authOK, err := createAuth(op.ArkAuth)
	if err != nil {
		return err
	}
	if !authOK {
		return fmt.Errorf("missing required field: sig")
	}
	parts := strings.Split(arkauth, ":") // fetch the signature from the string
	op.MsgClaimContractIncome.Signature, err = hex.DecodeString(parts[len(parts)-1])
	if err != nil {
		return fmt.Errorf("unable to decode signature: %s", err)
	}

	signer := sdk.MustAccAddressFromBech32(op.Signer)
	return sendMsg(&op.MsgClaimContractIncome, signer, op.Sequence, op, logs)
}

////////////////////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////////////////////

func sendMsg(msg sdk.Msg, signer sdk.AccAddress, seq *int64, op any, logs chan string) error {
	// check that message is valid
	err := msg.ValidateBasic()
	if err != nil {
		enc := json.NewEncoder(os.Stdout) // json instead of yaml to encode amount
		enc.SetIndent("", "  ")
		_ = enc.Encode(op)
		log.Fatal().Err(err).Msg("failed to validate basic")
	}

	// custom client context
	buf := bytes.NewBuffer(nil)
	ctx := clientCtx.WithFromAddress(signer)
	ctx = ctx.WithFromName(addressToName[signer.String()])
	ctx = ctx.WithOutput(buf)

	// override the sequence if provided
	txf := txFactory
	if seq != nil {
		txf = txFactory.WithSequence(uint64(*seq))
	}

	// send message
	err = tx.GenerateOrBroadcastTxWithFactory(ctx, txf, msg)
	if err != nil {
		fmt.Println(ColorPurple + "\nOperation:" + ColorReset)
		enc := json.NewEncoder(os.Stdout) // json instead of yaml to encode amount
		enc.SetIndent("", "  ")
		_ = enc.Encode(op)
		fmt.Println(ColorPurple + "\nTx Output:" + ColorReset)
		drainLogs(logs)
		return err
	}

	// extract txhash from output json
	var txRes sdk.TxResponse
	err = encodingConfig.Marshaler.UnmarshalJSON(buf.Bytes(), &txRes)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal tx response")
	}

	// fail if tx did not send, otherwise add to out native tx ids
	if txRes.Code != 0 {
		log.Debug().Uint32("code", txRes.Code).Str("log", txRes.RawLog).Msg("tx send failed")
	} else {
		nativeTxIDs = append(nativeTxIDs, txRes.TxHash)
	}

	return err
}
