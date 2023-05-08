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

### Features

- [#33](https://github.com/KYVENetwork/chain/pull/33) Upgrade Cosmos SDK to [v0.47.2](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.47.2) ([`v0.47.2-kyve-rc0`](https://github.com/KYVENetwork/cosmos-sdk/releases/tag/v0.47.2-kyve-rc0)).

## [v1.1.1](https://github.com/KYVENetwork/chain/releases/tag/v1.1.1) - 2023-05-08

### Improvements

- [#34](https://github.com/KYVENetwork/chain/pull/34) Support [Heighliner](https://github.com/strangelove-ventures/heighliner) to enable [interchaintest](https://github.com/strangelove-ventures/interchaintest).

## [v1.1.0](https://github.com/KYVENetwork/chain/releases/tag/v1.1.0) - 2023-04-18

### Improvements

- [#22](https://github.com/KYVENetwork/chain/pull/22) Various minor code improvements, cleanups, and validations.
- (deps) [#21](https://github.com/KYVENetwork/chain/pull/21), [#28](https://github.com/KYVENetwork/chain/pull/28) Bump Cosmos SDK to [v0.46.12](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.46.12) ([`v0.46.12-kyve`](https://github.com/KYVENetwork/cosmos-sdk/releases/tag/v0.46.12-kyve)).
- (deps) [#21](https://github.com/KYVENetwork/chain/pull/21) Switch to CometBFT from Informal Systems' Tendermint fork.
- (ibc) [#27](https://github.com/KYVENetwork/chain/pull/27) Enable tokens to be sent and received via IBC.

### Bug Fixes

- [#20](https://github.com/KYVENetwork/chain/pull/20) Adjust investor vesting schedules from second funding round.

### Client Breaking

- (`x/query`) [#23](https://github.com/KYVENetwork/chain/pull/23) Update the `StakerMetadata` query to reflect the new `Identity` and metadata fields.
- (`x/stakers`) [#23](https://github.com/KYVENetwork/chain/pull/23) Update `MsgUpdateMetadata` to reflect the new `Identity` and metadata fields.

### API Breaking

- [#22](https://github.com/KYVENetwork/chain/pull/22) Emit an event when updating module parameters.
- (`x/delegation`) [#24](https://github.com/KYVENetwork/chain/pull/24) Emit an event when a user initiates a protocol unbonding.
- (`x/pool`) [#24](https://github.com/KYVENetwork/chain/pull/24) Emit events for all module governance actions.
- (`x/stakers`) [#23](https://github.com/KYVENetwork/chain/pull/23) Update the event emitted when updating protocol node metadata.

### State Machine Breaking

- (`x/bundles`) [#19](https://github.com/KYVENetwork/chain/pull/19) Migrate `NetworkFee` param to type `sdk.Dec`.
- (`x/bundles`) [#22](https://github.com/KYVENetwork/chain/pull/22) Switch to a non-manipulable pseudo-random source seed for uploader selection.
- (`x/bundles`) [#26](https://github.com/KYVENetwork/chain/pull/26) Include the timestamp of the block that finalized a bundle.
- (`x/delegation`) [#19](https://github.com/KYVENetwork/chain/pull/19) Migrate `VoteSlash`, `UploadSlash`, `TimeoutSlash` params to type `sdk.Dec`.
- (`x/stakers`) [#19](https://github.com/KYVENetwork/chain/pull/19) Migrate `Commission` to type `sdk.Dec`.
- (`x/stakers`) [#23](https://github.com/KYVENetwork/chain/pull/23) Improve metadata by adding `Identity`, `SecurityContact`, `Details` fields, deprecating `Logo`.

## [v1.0.1](https://github.com/KYVENetwork/chain/releases/tag/v1.0.1) - 2023-05-08

### Improvements

- [#34](https://github.com/KYVENetwork/chain/pull/34) Support [Heighliner](https://github.com/strangelove-ventures/heighliner) to enable [interchaintest](https://github.com/strangelove-ventures/interchaintest).

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
