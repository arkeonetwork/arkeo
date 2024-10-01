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
	FoundationDevAccount       = "tarkeo1x978nttd8vgcgnv9wxut4dh7809lr0n2fhuh0q"
	FoundationCommunityAccount = "tarkeo124qmjmg55v6q5c5vy0vcpefrywxnxhkm7426pc"
	FoundationGrantsAccount    = "tarkeo1a307z4a82mcyv9njdj9ajnd9xpp90kmeqwntxj"
)
