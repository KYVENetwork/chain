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

- (deps) [#3](https://github.com/KYVENetwork/chain/pull/3) Bump Cosmos SDK to [v0.46.9](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.46.9) ([`v0.46.9-kyve-rc0`](https://github.com/KYVENetwork/cosmos-sdk/releases/tag/v0.46.9-kyve-rc0)).
- (deps) [#3](https://github.com/KYVENetwork/chain/pull/3) Bump IBC to [v6.1.0](https://github.com/cosmos/ibc-go/releases/tag/v6.1.0).

### State Machine Breaking

- (`x/bundles`) [#1](https://github.com/KYVENetwork/chain/pull/1) Migrate `StorageCost` param to type `sdk.Dec`.

## [v1.0.0-rc0](https://github.com/KYVENetwork/chain/releases/tag/v1.0.0-rc0) - 2023-02-03

`v1.0.0` Release Candidate for the Kaon network launch.
