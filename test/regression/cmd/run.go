package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// Run
////////////////////////////////////////////////////////////////////////////////////////

/*
type daemonLogger struct {
	name   string
	writer io.Writer
}

func (logger daemonLogger) Write(p []byte) (n int, err error) {
	padded := fmt.Sprintf("%-10s", fmt.Sprintf("[<%s>]", logger.name))
	return logger.writer.Write(append([]byte(padded), p...))
}

func newDaemonLogger(name string, writer io.Writer) daemonLogger {
	return daemonLogger{
		name:   name,
		writer: writer,
	}
}
*/

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

	providers := make([]types.Provider, 0)

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

		if op["type"] == "state" {
			// if we are pre-populating with providers, pull them out to insert
			// them into postgres as well
			if genesis, ok := op["genesis"]; ok {
				if app_state, ok := genesis.(map[string]any)["app_state"]; ok {
					if arkeo, ok := app_state.(map[string]any)["arkeo"]; ok {
						if provs, ok := arkeo.(map[string]any)["providers"]; ok {
							provsConvert, ok := provs.([]interface{})
							if ok {
								for _, p := range provsConvert {
									pk, ok := p.(map[string]any)["pub_key"].(string)
									if !ok {
										continue
									}
									ser, ok := p.(map[string]any)["service"].(int)
									if !ok {
										continue
									}
									provider := types.NewProvider(common.PubKey(pk), common.Service(ser))
									providers = append(providers, provider)
								}
							}
						}
					}
				}
			}
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

			// if its a state OpState and has provider, we need to push them into postgres
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

	// wait for postgres
	waitForPort("postgres", "directory-postgres:5432")
	iter := time.Now().UnixNano()
	tern(iter, providers) // create postgres db

	logLevel := "info"
	switch os.Getenv("DEBUG") {
	case "trace", "debug", "info", "warn", "error", "fatal", "panic":
		logLevel = os.Getenv("DEBUG")
	}

	sharedDirectoryEnv := []string{
		"DB_HOST=directory-postgres",
		"DB_PORT=5432",
		"DB_USER=arkeo",
		"DB_PASS=arkeo123",
		fmt.Sprintf("DB_NAME=arkeo_directory%d", iter),
		"DB_POOL_MAX_CONNS=2",
		"DB_POOL_MIN_CONNS=1",
		"DB_SSL_MODE=prefer",
	}
	procs := []process{
		{
			name:  "arkeod",
			cmd:   []string{"/regtest/cover-arkeod", "--log_level", logLevel, "start"},
			ports: []string{"8080", "26657"},
			env: []string{
				"GOCOVERDIR=/mnt/coverage",
			},
			sigkill: syscall.SIGKILL,
		},
		{
			name:  "sentinel",
			cmd:   []string{"/regtest/cover-sentinel"},
			ports: []string{"3636"},
			env: []string{
				"GOCOVERDIR=/mnt/coverage",
				fmt.Sprintf("PROVIDER_PUBKEY=%s", templatePubKey["pubkey_fox"]), // fox pubkey
				"NET=regtest",
				"MONIKER=regtest",
				"WEBSITE=n/a",
				"DESCRIPTION=n/a",
				"LOCATION=n/a",
				"PORT=3636",
				"SOURCE_CHAIN=http://localhost:1317",
				"EVENT_STREAM_HOST=localhost:26657",
				"FREE_RATE_LIMIT=10",
				"CLAIM_STORE_LOCATION=/regtest/.arkeo/claims",
				"CONTRACT_CONFIG_STORE_LOCATION=/regtest/.arkeo/contract_configs",
			},
			sigkill: syscall.SIGKILL,
		},
		{
			name:    "mock",
			cmd:     []string{"mock-daemon"},
			ports:   []string{"3765"},
			sigkill: syscall.SIGKILL,
		},
		{
			name:  "directory-api",
			cmd:   []string{"api"},
			ports: []string{"7777"},
			env: append(
				[]string{
					"LISTEN_ADDR=0.0.0.0:7777",
					"STATIC_DIR=/var/www/html",
				},
				sharedDirectoryEnv...,
			),
			sigkill: syscall.SIGKILL,
		},
		{
			name:  "directory-indexer",
			cmd:   []string{"indexer"},
			ports: []string{},
			env: append(
				[]string{
					"CHAIN_ID=arkeo",
					"BECH32_PREF_ACC_ADDR=tarkeo",
					"BECH32_PREF_ACC_PUB=tarkeopub",
					"ARKEO_API=http://localhost:1317",
					"TENDERMINT_API=http://localhost:26657",
					"TENDERMINT_WS=tcp://localhost:26657",
				},
				sharedDirectoryEnv...,
			),
			sigkill: syscall.SIGKILL,
		},
	}

	stderrLines := make(chan string, 100)
	for i, proc := range procs {
		if strings.HasPrefix(proc.name, "directory-") {
			for j := range proc.env {
				if strings.HasPrefix(proc.env[j], "DB_NAME=") {
					proc.env[j] = fmt.Sprintf("DB_NAME=arkeo_directory%d", iter)
				}
			}
		}
		procs[i].process = runProcess(proc, stderrLines)
	}

	// run the operations
	var returnErr error
	log.Info().Msgf("Executing %d operations", len(ops))
	for i, op := range ops {
		log.Info().Int("line", opLines[i]).Msgf(">>> [%d] %s", stateOpCount+i+1, op.OpType())
		returnErr = op.Execute(procs[0].process.Process, stderrLines)
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

	for _, proc := range procs {
		stopProcess(proc)
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

		iter = time.Now().UnixNano()
		tern(iter, providers)
		for _, proc := range procs {
			// restart process
			log.Debug().Msgf("Restarting process %s", proc.name)
			if strings.HasPrefix(proc.name, "directory-") {
				for j := range proc.env {
					if strings.HasPrefix(proc.env[j], "DB_NAME=") {
						proc.env[j] = fmt.Sprintf("DB_NAME=arkeo_directory%d", iter)
					}
				}
			}
			proc.process = runProcess(proc, stderrLines)
		}
	}

	return returnErr
}

type process struct {
	name    string
	cmd     []string
	ports   []string
	env     []string
	process *exec.Cmd
	sigkill syscall.Signal
}

func runProcess(proc process, stderrLines chan string) *exec.Cmd {
	// setup process io
	var process *exec.Cmd
	if len(proc.cmd) == 1 {
		process = exec.Command(proc.cmd[0], []string{}...) // #nosec G204
	} else if len(proc.cmd) > 1 {
		process = exec.Command(proc.cmd[0], proc.cmd[1:]...) // #nosec G204
	}
	process.Env = append(os.Environ(), proc.env...)
	stderr, err := process.StderrPipe()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to setup stderr process %s", proc.name)
	}
	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stderrScanner.Scan() {
			stderrLines <- fmt.Sprintf(">> %s > %s", proc.name, stderrScanner.Text())
		}
	}()
	if os.Getenv("DEBUG") != "" {
		process.Stdout = os.Stdout // newDaemonLogger(proc.name, os.Stdout)
		process.Stderr = os.Stderr // newDaemonLogger(proc.name, os.Stderr)
	}

	// start process
	log.Debug().Msgf("Starting process %s", proc.name)
	err = process.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to start process %s", proc.name)
	}

	// wait for process to listen on block creation port
	for _, port := range proc.ports {
		waitForPort(proc.name, fmt.Sprintf("localhost:%s", port))
	}
	return process
}

func stopProcess(proc process) {
	// stop process
	log.Debug().Msgf("Stopping process %s", proc.name)
	err := proc.process.Process.Signal(proc.sigkill)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to stop process %s", proc.name)
	}

	// wait for process to exit
	_, err = proc.process.Process.Wait()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to wait for process %s", proc.name)
	}
}

func tern(iter int64, providers []types.Provider) {
	// migrate postgres
	dbname := fmt.Sprintf("arkeo_directory%d", iter)
	log.Debug().Msgf("migrate postres %s", dbname)
	_, err := exec.Command("createdb", dbname).CombinedOutput()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to createdb")
	}
	cmd := exec.Command("tern", "migrate", "-c", "/app/directory/tern/tern.conf", "--database", dbname, "-m", "/app/directory/tern")
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("POSTGRES_DB=%s", dbname),
		"POSTGRES_USER=arkeo",
		"POSTGRES_PASSWORD=arkeo123",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Debug().Msg(string(out))
		log.Fatal().Err(err).Msg("failed to migrate postres")
	}

	// insert providers
	for _, provider := range providers {
		conf := db.DBConfig{
			Host:         "directory-postgres",
			Port:         5432,
			User:         "arkeo",
			Pass:         "arkeo123",
			DBName:       dbname,
			PoolMinConns: 0,
			PoolMaxConns: 5,
			SSLMode:      "prefer",
		}
		database, err := db.New(conf)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load db")
		}
		prov := db.ArkeoProvider{
			Pubkey:  provider.PubKey.String(),
			Service: provider.Service.String(),
			Bond:    "0",
		}
		_, err = database.InsertProvider(context.Background(), &prov)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to save provider")
		}
	}
}

func waitForPort(name, host string) {
	// wait for process to listen on block creation port
	log.Debug().Msgf("Waiting for %s port %s", name, host)
	for i := 0; ; i++ {
		time.Sleep(100 * time.Millisecond)
		conn, err := net.Dial("tcp", host)
		if err == nil {
			conn.Close()
			break
		}
		if i%100 == 0 {
			log.Debug().Msgf("Waiting for process to listen %s: %s", name, host)
		}
	}
}
