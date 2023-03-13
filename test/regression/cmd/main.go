package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////////////
// Main
////////////////////////////////////////////////////////////////////////////////////////

func main() {
	// parse the regex in the RUN environment variable to determine which tests to run
	runRegex := regexp.MustCompile(".*")
	if len(os.Getenv("RUN")) > 0 {
		runRegex = regexp.MustCompile(os.Getenv("RUN"))
	}

	// find all regression tests in path
	files := []string{}
	err := filepath.Walk("suites", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// skip files that are not yaml
		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		if runRegex.MatchString(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to find regression tests")
	}

	// sort the files
	sort.Strings(files)

	succeeded := []string{}
	failed := []string{}

	// run tests
	for _, file := range files {
		fmt.Println()

		// run test
		err := run(file)
		if err != nil {
			failed = append(failed, file)
			continue
		}

		// check export state
		err = export(file)
		if err != nil {
			failed = append(failed, file)
			continue
		}

		// success
		succeeded = append(succeeded, file)
	}

	// print the results
	fmt.Println()
	fmt.Printf("%sSucceeded:%s %d\n", ColorGreen, ColorReset, len(succeeded))
	for _, file := range succeeded {
		fmt.Printf("- %s\n", file)
	}
	fmt.Printf("%sFailed:%s %d\n", ColorRed, ColorReset, len(failed))
	for _, file := range failed {
		fmt.Printf("- %s\n", file)
	}
	fmt.Println()

	// exit with error code if any tests failed
	if len(failed) > 0 {
		os.Exit(1)
	}
}
