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

# v1.0.1-Prerelease 

## Added 
- Thorchain Claims Proto Updates
- Documentation of Testnet Setup using local build and Cosmovisor
- Documentation update and addition of validator setup documentation 
- New accounts to handle rewards 
- New Params to Arkeo Module 

## Changed
- Removed unused module account

## Fixed 
- Testnet binary generation using go build 
- Fixed Regression Testing 
- Updated Docs
- Fixed Consumer in Directory Service
- Fixed Regression Export 
- Fixed localnet docker 
- updated the genesis file
- claim timeout
- Fixed module imports
- update module to implement APPModuleBasic and AppModule
- Updated Tests on arkeo module keeper
- version issue on begin blocker 

# v1.0.0-Prerelease

## Added 

- Binary and Docker Image releaser 
- Release , Lint , Release Check into CI Actions
- Added ThorChain Claims and unit tests


## Changed 

- Updated Cosmos SDK from 0.46.13 to 0.50.8 
- Updated IBC 5 to 8.3.1
- Updated Proto generation using proto-builder image 

## Fixed

- Event Ordering 
- Updated Docker Dependencies 
- Fixed CI checks and Release Checks 
- Fixed UnitTest and Module Imports 
- Fixed Default Commands on the Arkeo Cmd