package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
)

// convertAddresses takes a list of addresses and converts them from oldPrefix to newPrefix.
func convertAddresses(addressList []string, oldPrefix, newPrefix string) ([]string, error) {
	var convertedAddresses []string

	for _, address := range addressList {
		hrp, data, err := bech32.Decode(address)
		if err != nil {
			return nil, fmt.Errorf("failed to decode address %s: %w", address, err)
		}

		// Check if the prefix matches the old prefix
		if hrp != oldPrefix {
			return nil, fmt.Errorf("expected prefix '%s', but got '%s' for address '%s'", oldPrefix, hrp, address)
		}

		// Encode with the new prefix
		mainnetAddress, err := bech32.Encode(newPrefix, data)
		if err != nil {
			return nil, fmt.Errorf("failed to encode address %s with new prefix %s: %w", address, newPrefix, err)
		}

		convertedAddresses = append(convertedAddresses, mainnetAddress)
	}
	return convertedAddresses, nil
}

// readAddressesFromFile reads a file line by line and returns a slice of addresses.
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

// writeAddressesToFile writes a slice of addresses to a file, one per line.
func writeAddressesToFile(filename string, addresses []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, address := range addresses {
		_, err := writer.WriteString(address + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", filename, err)
		}
	}

	return writer.Flush()
}

func main() {
	// Define flags for CLI
	oldPrefix := flag.String("oldPrefix", "", "Old prefix of the addresses")
	newPrefix := flag.String("newPrefix", "", "New prefix to replace the old prefix")
	addresses := flag.String("addresses", "", "Comma-separated list of addresses to convert")
	addressFile := flag.String("file", "", "File containing list of addresses to convert (one per line)")

	// Parse command-line arguments
	flag.Parse()

	// Check for required arguments
	if *oldPrefix == "" || *newPrefix == "" {
		log.Fatal("Usage: go run converter.go -oldPrefix=<oldPrefix> -newPrefix=<newPrefix> -addresses=<comma-separated addresses> or -file=<filename>")
	}

	var addressList []string
	var err error

	// Determine source of addresses (file or list)
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

	// Convert addresses
	convertedAddresses, err := convertAddresses(addressList, *oldPrefix, *newPrefix)
	if err != nil {
		log.Fatalf("Error converting addresses: %v", err)
	}

	// Write converted addresses to file
	outputFile := "converted.txt"
	err = writeAddressesToFile(outputFile, convertedAddresses)
	if err != nil {
		log.Fatalf("Error writing converted addresses to file: %v", err)
	}

	fmt.Printf("Converted addresses written to %s\n", outputFile)
}
