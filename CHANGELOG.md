<!--

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

-->

# CHANGELOG

## [Unreleased]

### Improvements

- (deps) [#21](https://github.com/KYVENetwork/chain/pull/21) Bump Cosmos SDK to [v0.46.11](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.46.11) ([`v0.46.11-kyve-rc0`](https://github.com/KYVENetwork/cosmos-sdk/releases/tag/v0.46.11-kyve-rc0)).
- (deps) [#21](https://github.com/KYVENetwork/chain/pull/21) Switch to CometBFT from Informal Systems' Tendermint fork.

### Bug Fixes

- [#20](https://github.com/KYVENetwork/chain/pull/20) Adjust investor vesting schedules from second funding round.

## [v1.0.0](https://github.com/KYVENetwork/chain/releases/tag/v1.0.0) - 2023-03-10

Release for the KYVE network launch.

## [v1.0.0-rc1](https://github.com/KYVENetwork/chain/releases/tag/v1.0.0-rc1) - 2023-03-07

`v1.0.0` Release Candidate for a Kaon network upgrade.

### Improvements

- (deps) [#3](https://github.com/KYVENetwork/chain/pull/3), [#7](https://github.com/KYVENetwork/chain/pull/7) Bump Cosmos SDK to [v0.46.10](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.46.10) ([`v0.46.10-kyve-rc0`](https://github.com/KYVENetwork/cosmos-sdk/releases/tag/v0.46.10-kyve-rc0)).
- (deps) [#3](https://github.com/KYVENetwork/chain/pull/3) Bump IBC to [v6.1.0](https://github.com/cosmos/ibc-go/releases/tag/v6.1.0).
- (deps) [#7](https://github.com/KYVENetwork/chain/pull/7) Bump Tendermint to [v0.34.26](https://github.com/informalsystems/tendermint/releases/tag/v0.34.26).
- (`x/team`) [#7](https://github.com/KYVENetwork/chain/pull/7) Switch to a co-minting approach.

### State Machine Breaking

- (`x/bundles`) [#1](https://github.com/KYVENetwork/chain/pull/1) Migrate `StorageCost` param to type `sdk.Dec`.

## [v1.0.0-rc0](https://github.com/KYVENetwork/chain/releases/tag/v1.0.0-rc0) - 2023-02-03

`v1.0.0` Release Candidate for the Kaon network launch.
