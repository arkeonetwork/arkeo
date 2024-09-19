package types

const (
	// ModuleName defines the module name
	ModuleName   = "arkeo"
	ReserveName  = "arkeo-reserve"
	ProviderName = "providers"
	ContractName = "contracts"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_arkeo"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// Foundational Accounts
const (
	FoundationDevAccount       = "tarkeo10sav33v67743s6cl2cvjmmua7c5arysw3txz9r"
	FoundationCommunityAccount = "tarkeo1v50hrsxx0mxar4653aujcnqyjft07w0npcxrjx"
	FoundationGrantsAccount    = "tarkeo16k3k0erkwaanqnup20dxxenpd6wh058nh4pgup"
)
