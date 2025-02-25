# Changelog

All notable changes to this project will be documented in this file.

<!--
### Added

Contains the new features.

### Changed

Contains API breaking changes to existing functionality.

### Deprecated

Contains the candidates for removal in a future release.

### Removed

Contains API breaking changes of removed APIs.

### Fixed

Contains bug fixes.

### Improvements

Contains all the PRs that improved the code without changing the behaviors.
-->

## v1.0.6-Prerelease

### Added
- Added `VersionForAddress` field to `GenesisState`.

### Changed
- Updated outdated dependencies.
- Moved funds from claim to liquidity reserve when airdrop ends.

### Fixed
- Fixed old expirations not being removed.
- Fixed incorrect delegate usage in `msg open contract`.
- Fixed the recipient’s `IsTransferable` field being overwritten to `false` in `MsgClaimThorchain`.
- Fixed `GenesisState` validation.
- Fixed validator rewards payout.
- Fixed bug allowing double claims with Thorchain claims.
- Fixed non deterministic map iteration to sorted iteration 
- Fixed fee stuck in arkeo module by moving fees to arkeo-reserve
- Fixed signature replay on different network
- Fixed informational issues
- Fixed claim validation

---

## v1.0.5-Prerelease

### Added
- Added Arkeo testnet validator addresses to airdrop.

### Changed
- Updated Docker images.
- Set minimum gas price to zero.

### Fixed
- Updated event labels on Thorchain delegate events.
- Added `Response` to claim functions.

---

## v1.0.4-Prerelease

### Changed
- Updated Docker images.

---

## v1.0.3-Prerelease

### Added
- Added Sentinel setup documentation.
- Added Sentinel regression test.
- Added GoReleaser for Sentinel and Directory Service.

### Changed
- Updated Sentinel to handle provider events.
- Separated testnet and mainnet releaser.
- Updated Thorchain claim address.

### Fixed
- Fixed code lint issues.
- Fixed WebSocket client issue with event stream.
- Fixed Swagger issue.
- Fixed Docker files and scripts.
- Fixed GoReleaser for binary and Docker.

---

## v1.0.2-Prerelease

### Changed
- Updated Thorchain claim server address handling.

### Fixed
- Fixed regression test version issues.

---

## v1.0.1-Prerelease

### Added
- Added claim record scripts.
- Updated Thorchain Claims Proto.
- Added testnet setup documentation using local build and Cosmovisor.
- Added validator setup documentation.
- Added new accounts to handle rewards.
- Introduced new parameters to Arkeo module.
- Updated validator reward distribution.
- Released testnet Genesis file.

### Changed
- Removed unused module account.
- Disabled system validator rewards.
- Set default mint parameters to zero.
- Redirected validator and delegator rewards from the Reserve Module.
- Updated testnet seed and peer addresses.
- Updated testnet documentation.

### Fixed
- Fixed testnet binary generation using `go build`.
- Fixed regression testing issues.
- Updated documentation.
- Fixed consumer in Directory Service.
- Fixed regression export.
- Fixed localnet Docker setup.
- Updated genesis file.
- Fixed claim timeout.
- Fixed module imports.
- Updated module to implement `APPModuleBasic` and `AppModule`.
- Updated tests for Arkeo module keeper.
- Fixed version issue in `begin_blocker`.
- Fixed Genesis URL.
- Fixed Thorchain claim server address on mainnet.

---

## v1.0.0-Prerelease

### Added
- Introduced binary and Docker image releaser.
- Implemented CI actions for release, lint, and release check.
- Added Thorchain claims and unit tests.

### Changed
- Upgraded Cosmos SDK from `0.46.13` to `0.50.8`.
- Upgraded IBC from `5` to `8.3.1`.
- Updated Proto generation using `proto-builder` image.

### Fixed
- Fixed event ordering issues.
- Updated Docker dependencies.
- Fixed CI checks and release checks.
- Fixed unit tests and module imports.
- Fixed default commands in the Arkeo CLI.

---

