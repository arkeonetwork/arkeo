package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"text/template"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

////////////////////////////////////////////////////////////////////////////////////////
// Run
////////////////////////////////////////////////////////////////////////////////////////

func run(path string) error {
	log.Info().Msgf("Running regression test: %s", path)

	// reset native txids
	nativeTxIDs = nativeTxIDs[:0]

	// clear data directory
	log.Debug().Msg("Clearing data directory")
	out, err := exec.Command("rm", "-rf", "/regtest/.arkeo").CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		log.Fatal().Err(err).Msg("failed to clear data directory")
	}

	// init chain with dog mnemonic
	log.Debug().Msg("Initializing chain")
	cmd := exec.Command("arkeod", "init", "local", "--chain-id", "arkeo", "--staking-bond-denom", "uarkeo", "--recover")
	cmd.Stdin = bytes.NewBufferString(dogMnemonic + "\n")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		log.Fatal().Err(err).Msg("failed to initialize chain")
	}

	// init chain
	log.Debug().Msg("Initializing chain")
	cmd = exec.Command("arkeod", "init", "local", "--chain-id", "arkeo", "--staking-bond-denom", "uarkeo", "-o")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		log.Fatal().Err(err).Msg("failed to initialize chain")
	}

	// clone common templates
	tmpls := template.Must(templates.Clone())

	// ensure no naming collisions
	if tmpls.Lookup(filepath.Base(path)) != nil {
		log.Fatal().Msgf("test name collision: %s", filepath.Base(path))
	}

	// read the file
	f, err := os.Open(path)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open test file")
	}
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read test file")
	}
	f.Close()

	// track line numbers
	opLines := []int{0}
	scanner := bufio.NewScanner(bytes.NewBuffer(fileBytes))
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if line == "---" {
			opLines = append(opLines, i+2)
		}
	}

	// parse the template
	tmpl, err := tmpls.Parse(string(fileBytes))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse template")
	}

	// render the template
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to render template")
	}

	// all operations we will execute
	ops := []Operation{}

	// track whether we've seen non-state operations
	seenNonState := false

	dec := yaml.NewDecoder(buf)
	for {
		// decode into temporary type
		op := map[string]any{}
		err = dec.Decode(&op)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal().Err(err).Msg("failed to decode operation")
		}

		// warn empty operations
		if len(op) == 0 {
			log.Warn().Msg("empty operation")
			continue
		}

		// state operations must be first
		if op["type"] == "state" && seenNonState {
			log.Fatal().Msg("state operations must be first")
		}
		if op["type"] != "state" {
			seenNonState = true
		}

		ops = append(ops, NewOperation(op))
	}

	// warn if no operations found
	if len(ops) == 0 {
		err = errors.New("no operations found")
		log.Err(err).Msg("")
		return err
	}

	// execute all state operations
	stateOpCount := 0
	for i, op := range ops {
		if _, ok := op.(*OpState); ok {
			log.Info().Int("line", opLines[i]).Msgf(">>> [%d] %s", i+1, op.OpType())
			err = op.Execute(nil, nil)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to execute state operation")
			}
			stateOpCount++
		}
	}
	ops = ops[stateOpCount:]
	opLines = opLines[stateOpCount:]

	// validate genesis
	log.Debug().Msg("Validating genesis")
	cmd = exec.Command("arkeod", "validate-genesis")
	out, err = cmd.CombinedOutput()
	if err != nil {
		// dump the genesis
		fmt.Println(ColorPurple + "Genesis:" + ColorReset)
		f, err := os.OpenFile("/regtest/.arkeo/config/genesis.json", os.O_RDWR, 0o644)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to open genesis file")
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		f.Close()

		// dump error and exit
		fmt.Println(string(out))
		log.Fatal().Err(err).Msg("genesis validation failed")
	}

	// overwrite private validator key
	log.Debug().Msg("Overwriting private validator key")
	cmd = exec.Command("cp", "/mnt/priv_validator_key.json", "/regtest/.arkeo/config/priv_validator_key.json")
	err = cmd.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to overwrite private validator key")
	}

	// overwrite tendermint config
	log.Debug().Msg("Overwriting tendermint config")
	cmd = exec.Command("cp", "/mnt/config.toml", "/regtest/.arkeo/config/config.toml")
	err = cmd.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to overwrite private validator key")
	}

	// overwrite cosmos config
	log.Debug().Msg("Overwriting tendermint config")
	cmd = exec.Command("cp", "/mnt/app.toml", "/regtest/.arkeo/config/app.toml")
	err = cmd.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to overwrite private validator key")
	}

	// setup process io
	arkeo := exec.Command("/regtest/cover-arkeod", "start")
	arkeo.Env = append(os.Environ(), "GOCOVERDIR=/mnt/coverage")
	stderr, err := arkeo.StderrPipe()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup arkeo stderr")
	}
	stderrScanner := bufio.NewScanner(stderr)
	stderrLines := make(chan string, 100)
	go func() {
		for stderrScanner.Scan() {
			stderrLines <- stderrScanner.Text()
		}
	}()
	if os.Getenv("DEBUG") != "" {
		arkeo.Stdout = os.Stdout
		arkeo.Stderr = os.Stderr
	}

	// start arkeo process
	log.Debug().Msg("Starting arkeod")
	err = arkeo.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start arkeod")
	}

	// wait for arkeo to listen on block creation port
	for i := 0; ; i++ {
		time.Sleep(100 * time.Millisecond)
		conn, err := net.Dial("tcp", "localhost:8080")
		if err == nil {
			conn.Close()
			break
		}
		if i%100 == 0 {
			log.Debug().Msg("Waiting for arkeo to listen")
		}
	}

	// wait for arkeo to listen on block rpc port
	for i := 0; ; i++ {
		time.Sleep(100 * time.Millisecond)
		conn, err := net.Dial("tcp", "localhost:26657")
		if err == nil {
			conn.Close()
			break
		}
		if i%100 == 0 {
			log.Debug().Msg("Waiting for arkeo to listen")
		}
	}

	// setup process for sentinel
	sentinel := exec.Command("/regtest/cover-sentinel")
	sentinel.Env = append(
		os.Environ(),
		"GOCOVERDIR=/mnt/coverage",
		fmt.Sprintf("PROVIDER_PUBKEY=%s", templatePubKey["pubkey_fox"]), // fox pubkey
		"NET=regtest",
		"MONIKER=regtest",
		"WEBSITE=n/a",
		"DESCRIPTION=n/a",
		"LOCATION=n/a",
		"PORT=3636",
		"PROXY_HOST=https://swapi.dev", // TODO: remove me
		"SOURCE_CHAIN=localhost:1317",
		"EVENT_STREAM_HOST=localhost:26657",
		"FREE_RATE_LIMIT=10",
		"FREE_RATE_LIMIT_DURATION=1m",
		"SUB_RATE_LIMIT=10",
		"SUB_RATE_LIMIT_DURATION=1m",
		"AS_GO_RATE_LIMIT=10",
		"AS_GO_RATE_LIMIT_DURATION=1m",
		"CLAIM_STORE_LOCATION=/regtest/.arkeo/claims",
		"GAIA_RPC_ARCHIVE_HOST=http://176.34.207.130:26657",
	)
	stderr, err = sentinel.StderrPipe()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup sentinel stderr")
	}
	stderrScanner2 := bufio.NewScanner(stderr)
	stderrLines2 := make(chan string, 100)
	go func() {
		for stderrScanner2.Scan() {
			stderrLines2 <- stderrScanner.Text()
		}
	}()
	if os.Getenv("DEBUG") != "" {
		sentinel.Stdout = os.Stdout
		sentinel.Stderr = os.Stderr
	}

	// start arkeo process
	log.Debug().Msg("Starting sentinel")
	err = sentinel.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start sentinel")
	}

	// wait for sentinel to listen on block creation port
	for i := 0; ; i++ {
		time.Sleep(100 * time.Millisecond)
		conn, err := net.Dial("tcp", "localhost:3636")
		if err == nil {
			conn.Close()
			break
		}
		if i%100 == 0 {
			log.Debug().Msg("Waiting for sentinel to listen")
		}
	}

	// run the operations
	var returnErr error
	log.Info().Msgf("Executing %d operations", len(ops))
	for i, op := range ops {
		log.Info().Int("line", opLines[i]).Msgf(">>> [%d] %s", stateOpCount+i+1, op.OpType())
		returnErr = op.Execute(arkeo.Process, stderrLines)
		if returnErr != nil {
			log.Error().Err(returnErr).
				Int("line", opLines[i]).
				Int("op", stateOpCount+i+1).
				Str("type", op.OpType()).
				Str("path", path).
				Msg("operation failed")
			fmt.Println()
			dumpLogs(stderrLines)
			break
		}
	}

	// log success
	if returnErr == nil {
		log.Info().Msg("All operations succeeded")
	}

	// stop arkeo process
	log.Debug().Msg("Stopping arkeod")
	err = arkeo.Process.Signal(syscall.SIGUSR1)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to stop arkeod")
	}

	// wait for process to exit
	_, err = arkeo.Process.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to wait for arkeod")
	}

	// stop sentinel process
	log.Debug().Msg("Stopping sentinel")
	err = sentinel.Process.Signal(syscall.SIGKILL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to stop sentinel")
	}

	// wait for process to exit
	_, err = sentinel.Process.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to wait for sentinel")
	}

	// if failed and debug enabled restart to allow inspection
	if returnErr != nil && os.Getenv("DEBUG") != "" {

		// remove validator key (otherwise arkeo will hang in begin block)
		log.Debug().Msg("Removing validator key")
		cmd = exec.Command("rm", "/regtest/.arkeo/config/priv_validator_key.json")
		out, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			log.Fatal().Err(err).Msg("failed to remove validator key")
		}

		// restart arkeo
		log.Debug().Msg("Restarting arkeod")
		arkeo = exec.Command("arkeod", "start")
		arkeo.Stdout = os.Stdout
		arkeo.Stderr = os.Stderr
		err = arkeo.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to restart arkeod")
		}

		// restart sentinel
		log.Debug().Msg("Restarting sentinel")
		sentinel = exec.Command("sentinel")
		sentinel.Stdout = os.Stdout
		sentinel.Stderr = os.Stderr
		err = sentinel.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to restart sentinel")
		}

		// wait for arkeo
		log.Debug().Msg("Waiting for arkeod")
		_, err = arkeo.Process.Wait()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to wait for arkeod")
		}

		// wait for sentinel
		log.Debug().Msg("Waiting for sentinel")
		_, err = sentinel.Process.Wait()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to wait for sentinel")
		}
	}

	return returnErr
}
