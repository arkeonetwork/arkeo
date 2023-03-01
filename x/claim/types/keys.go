package types

const (
	// ModuleName defines the module name
	ModuleName = "arkeo-claim" // preappend arkeo to avoid namespace collisions with comsos claim module

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_claim"

	// ClaimRecordsStorePrefix defines the store prefix for the claim records (by arkeo address)
	ClaimRecordsArkeoStorePrefix = "claimrecordsarkeo"

	// ClaimRecordsStorePrefix defines the store prefix for the claim records (by eth address)
	ClaimRecordsEthStorePrefix = "claimrecordsethereum"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
