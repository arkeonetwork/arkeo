package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func deepMerge(a, b map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		switch vv := v.(type) {
		case []any:
			if bv, ok := result[k]; ok {
				if bv, ok := bv.([]any); ok {
					result[k] = append(bv, vv...)
					continue
				}
			}
		case map[string]any:
			if bv, ok := result[k]; ok {
				if bv, ok := bv.(map[string]any); ok {
					result[k] = deepMerge(bv, vv)
					continue
				}
			}
		}
		result[k] = v
	}
	return result
}

func dumpLogs(logs chan string) {
	for {
		select {
		case line := <-logs:
			fmt.Println(ColorPurple + ">>> " + ColorReset + line)
			continue
		default:
		}
		break
	}
}

func drainLogs(logs chan string) {
	// if DEBUG is set skip draining logs
	if os.Getenv("DEBUG") != "" {
		return
	}

	for {
		select {
		case <-logs:
			continue
		case <-time.After(100 * time.Millisecond):
		}
		break
	}
}

func processRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(os.Signal(nil))
	return err == nil
}

func getTimeFactor() time.Duration {
	tf, err := strconv.ParseInt(os.Getenv("TIME_FACTOR"), 10, 64)
	if err != nil {
		return time.Duration(1)
	}
	return time.Duration(tf)
}
