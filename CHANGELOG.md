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

# v1.0.3-Prerelease 

## Added 
- Added sentinel setup docs

## Changed
- Updated sentinel to handle provider events

## Fixed 
- Fixed code lint
- Fixed ws client issue with event stream
- Fixed swagger issue

# v1.0.2-Prerelease 

## Added 

## Changed
- Updated thorchain claim server address handling

## Fixed 
- Fixed Regression Test Version Issues 

# v1.0.1-Prerelease 

## Added 
- Thorchain Claims Proto Updates
- Documentation of Testnet Setup using local build and Cosmovisor
- Documentation update and addition of validator setup documentation 
- New accounts to handle rewards 
- New Params to Arkeo Module 
- Update validator distribution of rewards 
- Testnet Genesis File

## Changed
- Removed unused module account
- Disabled System Validator Rewards 
- Default Mint params set to zero
- Validator and Delegator rewards from Reserve Module
- Updated testnet seed and peer address
- Updated testnet docs 

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
- Fixed Genesis Url
- Fixed thorchain claim server address mainnet

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