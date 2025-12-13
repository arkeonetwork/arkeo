package types

// Hardcoded extra authorities allowed to mutate the service registry.
// TODO: replace with on-chain parameter or governance-driven allowlist.
var ExtraAuthorities = map[string]struct{}{
	"arkeo1w2ln0prejgrztmf9w23e0rsnlks7djneh5te7p": {},
}

// IsAuthorityAllowed returns true if addr matches the keeper authority or is in ExtraAuthorities.
func IsAuthorityAllowed(keeperAuthority, addr string) bool {
	if addr == keeperAuthority {
		return true
	}
	if _, ok := ExtraAuthorities[addr]; ok {
		return true
	}
	return false
}
