package main

import (
	"fmt"
	"text/template"
	"time"
)

////////////////////////////////////////////////////////////////////////////////////////
// Templates
////////////////////////////////////////////////////////////////////////////////////////

// nativeTxIDs will be reset on each run and contains the native txids for all sent txs
var nativeTxIDs = []string{}

// templates contain all base templates referenced in tests
var templates *template.Template

// funcMap is a map of functions that can be used in all templates and tests
var funcMap = template.FuncMap{
	"observe_txid": func(i int) string {
		return fmt.Sprintf("%064x", i) // padded 64-bit hex string
	},
	"native_txid": func(i int) string {
		// this will get double-rendered
		if len(nativeTxIDs) == 0 {
			return fmt.Sprintf("{{ native_txid %d }}", i)
		}
		// allow reverse indexing
		if i < 0 {
			i += len(nativeTxIDs) + 1
		}
		return nativeTxIDs[i-1]
	},
	"addr_module_bonded_tokens_pool": func() string {
		return ModuleAddrBondedTokensPool
	},
	"addr_module_not_bonded_tokens_pool": func() string {
		return ModuleAddrNotBondedTokensPool
	},
	"addr_module_gov": func() string {
		return ModuleAddrGov
	},
	"addr_module_distribution": func() string {
		return ModuleAddrDistribution
	},
	"addr_module_fee_collector": func() string {
		return ModuleAddrFeeCollector
	},
	"addr_module_arkeo": func() string {
		return ModuleAddr
	},
	"addr_module_provider": func() string {
		return ModuleAddrProvider
	},
	"addr_module_contract": func() string {
		return ModuleAddrContract
	},
	"addr_module_reserve": func() string {
		return ModuleAddrReserve
	},
	"timestamp": func() int64 {
		return time.Now().UnixNano()
	},
}

////////////////////////////////////////////////////////////////////////////////////////
// Functions
////////////////////////////////////////////////////////////////////////////////////////

func init() {
	// register template names for all keys
	for k, v := range templateAddress {
		vv := v // copy
		funcMap[k] = func() string {
			return vv
		}
	}
	for k, v := range templatePubKey {
		vv := v // copy
		funcMap[k] = func() string {
			return vv
		}
	}

	// parse all templates with custom functions
	templates = template.Must(
		template.New("").Funcs(funcMap).ParseGlob("templates/*.yaml"),
	)
}
