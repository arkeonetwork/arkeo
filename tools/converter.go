package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/types"
)

// GetModuleAddress computes the address for a given module name with a specific prefix.
func GetModuleAddress(moduleName string) types.AccAddress {
	return types.AccAddress(crypto.AddressHash([]byte(moduleName)))
}

// convertAddresses takes a list of addresses and converts them from oldPrefix to newPrefix.
func convertAddresses(addressList []string, oldPrefix, newPrefix string) (map[string]string, error) {
	convertedMap := make(map[string]string)

	for _, address := range addressList {
		hrp, data, err := bech32.Decode(address)
		if err != nil {
			return nil, fmt.Errorf("failed to decode address %s: %w", address, err)
		}

		// Check if the prefix matches the old prefix
		if hrp != oldPrefix {
			return nil, fmt.Errorf("expected prefix '%s', but got '%s' for address '%s'", oldPrefix, hrp, address)
		}

		mainnetAddress, err := bech32.Encode(newPrefix, data)
		if err != nil {
			return nil, fmt.Errorf("failed to encode address %s with new prefix %s: %w", address, newPrefix, err)
		}

		convertedMap[address] = mainnetAddress
	}
	return convertedMap, nil
}

func readAddressesFromFile(filename string) ([]string, error) {
	var addresses []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		addresses = append(addresses, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read from file %s: %w", filename, err)
	}

	return addresses, nil
}

// writeAddressesToJSONFile writes a map of old-to-new addresses to a JSON file.
func writeAddressesToJSONFile(filename string, addressMap map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(addressMap); err != nil {
		return fmt.Errorf("failed to write JSON to file %s: %w", filename, err)
	}

	return nil
}

func main() {
	// Define flags for CLI
	oldPrefix := flag.String("oldPrefix", "", "Old prefix of the addresses")
	newPrefix := flag.String("newPrefix", "", "New prefix to replace the old prefix")
	addresses := flag.String("addresses", "", "Comma-separated list of addresses to convert")
	addressFile := flag.String("file", "", "File containing list of addresses to convert (one per line)")
	moduleName := flag.String("module", "", "Module name to generate the address for")

	// Parse command-line arguments
	flag.Parse()

	cfg := types.GetConfig()

	if *newPrefix == "arkeo" {
		cfg.SetBech32PrefixForAccount("arkeo", "arkeopub")
	} else if *newPrefix == "tarkeo" {
		cfg.SetBech32PrefixForAccount("tarkeo", "tarkeopub")
	} else {
		log.Fatalf("Unsupported module prefix '%s'. Supported prefixes are 'arkeo' and 'tarkeo'.", *newPrefix)
	}

	// Module address generation operation
	if *moduleName != "" {

		moduleAddr := GetModuleAddress(*moduleName)

		fmt.Printf("Generated address for module '%s': %s\n", *moduleName, moduleAddr)
		return
	}

	// Address conversion operation
	if *oldPrefix == "" || *newPrefix == "" {
		log.Fatal("Usage: go run converter.go -oldPrefix=<oldPrefix> -newPrefix=<newPrefix> -addresses=<comma-separated addresses> or -file=<filename>")
	}

	var addressList []string
	var err error

	if *addressFile != "" {
		addressList, err = readAddressesFromFile(*addressFile)
		if err != nil {
			log.Fatalf("Error reading addresses from file: %v", err)
		}
	} else if *addresses != "" {
		addressList = strings.Split(*addresses, ",")
	} else {
		log.Fatal("Please provide addresses either through -addresses or -file option.")
	}

	convertedMap, err := convertAddresses(addressList, *oldPrefix, *newPrefix)
	if err != nil {
		log.Fatalf("Error converting addresses: %v", err)
	}

	outputFile := "converted.json"
	err = writeAddressesToJSONFile(outputFile, convertedMap)
	if err != nil {
		log.Fatalf("Error writing converted addresses to file: %v", err)
	}

	fmt.Printf("Converted addresses written to %s", outputFile)
}
