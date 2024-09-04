package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog/log"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// Export
////////////////////////////////////////////////////////////////////////////////////////

func export(path string) error {
	// export state
	log.Debug().Msg("Exporting state")
	cmd := exec.Command("arkeod", "export")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to export state")
	}

	// decode export
	var export map[string]any
	err = json.Unmarshal(out, &export)
	if err != nil {
		fmt.Println(string(out))
		log.Fatal().Err(err).Msg("failed to decode export")
	}

	// ignore genesis time and version for comparison
	delete(export, "genesis_time")

	// ignore any version or time fields in app state
	appState, _ := export["app_state"].(map[string]any)
	claimArkeo, _ := appState["claimarkeo"].(map[string]any)
	params, _ := claimArkeo["params"].(map[string]any)
	delete(params, "airdrop_start_time")
	staking, _ := appState["staking"].(map[string]any)
	validators, _ := staking["validators"].([]any)
	for i, validator := range validators {
		v, _ := validator.(map[string]any)
		commission, _ := v["commission"].(map[string]any)
		delete(commission, "update_time")
		v["commission"] = commission
		validators[i] = v
	}

	// encode export
	out, err = json.MarshalIndent(export, "", "  ")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to encode export")
	}

	// base path without extension and replace path separators with underscores
	exportName := strings.TrimSuffix(path, filepath.Ext(path))
	exportName = strings.ReplaceAll(exportName, string(os.PathSeparator), "_")
	exportPath := fmt.Sprintf("/mnt/exports/%s.json", exportName)

	// check whether existing export exists
	_, err = os.Stat(exportPath)
	exportExists := err == nil

	// check export invariants
	err = checkExportInvariants(export)
	if err != nil {
		// also log export changes for easier debugging
		if exportExists {
			_ = checkExportChanges(export, exportPath)
		}

		return err
	}

	// export if it none exists or EXPORT is set
	if !exportExists || os.Getenv("EXPORT") != "" {
		log.Debug().Msg("Writing export")
		err = os.WriteFile(exportPath, out, 0o600)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to write export")
		}
		return nil
	}

	return checkExportChanges(export, exportPath)
}

////////////////////////////////////////////////////////////////////////////////////////
// Checks
////////////////////////////////////////////////////////////////////////////////////////

func checkExportInvariants(genesis map[string]any) error {
	// check export invariants
	log.Debug().Msg("Checking export invariants")
	appState, _ := genesis["app_state"].(map[string]any)

	// encode arkeonetwork state to json for custom unmarshal
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	err := enc.Encode(appState["arkeo"])
	if err != nil {
		log.Fatal().Err(err).Msg("failed to encode genesis state")
	}

	// unmarshal json to genesis state
	genesisState := &types.GenesisState{}
	err = encodingConfig.Marshaler.UnmarshalJSON(buf.Bytes(), genesisState)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to decode genesis state")
	}

	// TODO: add any arkeo specific invariant checks

	return err
}

func checkExportChanges(newExport map[string]any, path string) error {
	// compare existing export
	log.Debug().Msg("Reading existing export")

	// open existing export
	f, err := os.Open(path)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open existing export")
	}
	defer f.Close()

	// decode existing export
	oldExport := map[string]any{}
	err = json.NewDecoder(f).Decode(&oldExport)
	if err != nil {
		log.Err(err).Msg("failed to decode existing export")
	}

	// compare exports
	log.Debug().Msg("Comparing exports")
	diff := cmp.Diff(oldExport, newExport)
	if diff != "" {
		log.Error().Msgf("exports differ: %s", diff)
		return errors.New("exports differ")
	}

	log.Info().Msg("State export matches expected")
	return nil
}
