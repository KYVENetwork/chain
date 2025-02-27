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

An '!' indicates a state machine breaking change.

## [v2.0.0](https://github.com/KYVENetwork/chain/releases/tag/v2.0.0) - 2025-02-27

### Features

- ! (`x/bundles`) [#236](https://github.com/KYVENetwork/chain/pull/236) merkle_root bundles migration
- ! (`x/stakers`) [#209](https://github.com/KYVENetwork/chain/pull/209) Shared Staking: Consensus-validator stake is now used for the protocol.
- ! (`x/stakers`) [#210](https://github.com/KYVENetwork/chain/pull/210) Shared Staking: Pool specific commission and stake fraction.
- ! (`x/stakers`) [#211](https://github.com/KYVENetwork/chain/pull/211) Shared Staking: Maximum voting power per pool.

### Improvements

- (deps) [#250](https://github.com/KYVENetwork/chain/pull/250) Remove group module.
- (`x/stakers`) [#232](https://github.com/KYVENetwork/chain/pull/232) Shared Staking: Update stakers queries 
- (deps) [#213](https://github.com/KYVENetwork/chain/pull/213) Bump to CosmosSDK v0.50.11 and cleanup deps.
- (deps) [#246](https://github.com/KYVENetwork/chain/pull/246) Bump stakers module version from v1beta1 to v1

### Bug Fixes

- ! [#207](https://github.com/KYVENetwork/chain/pull/207) Added test case for points.

## [v1.5.0](https://github.com/KYVENetwork/chain/releases/tag/v1.5.0) - 2024-07-05

### Features

- ! (`x/stakers`) [#185](https://github.com/KYVENetwork/chain/pull/185) Support multiple coins in stakers module.
- ! (`x/bundles`) [#184](https://github.com/KYVENetwork/chain/pull/184) Support multiple coins in bundles module.
- ! (`x/funders`) [#179](https://github.com/KYVENetwork/chain/pull/179) Support multiple coins for funding.
- ! (`x/bundles`) [#177](https://github.com/KYVENetwork/chain/pull/177) Add end-key to pool.
- ! (`x/bundles`) [#178](https://github.com/KYVENetwork/chain/pull/178) Tally bundle proposal on timeout.
- ! (`x/delegation`) [#176](https://github.com/KYVENetwork/chain/pull/176) Support multiple denoms for delegation rewards.
- ! (`x/bundles`) [#172](https://github.com/KYVENetwork/chain/pull/172) Storage-cost per storage-provider.
- ! (`x/bundles`) [#175](https://github.com/KYVENetwork/chain/pull/175) Tally bundle proposal on skip-uploader-role.

### Improvements

- [#199](https://github.com/KYVENetwork/chain/pull/199) Align events naming.
- ! (deps) [#198](https://github.com/KYVENetwork/chain/pull/198) Bump ibc version to v8.3.1.
- (deps) [#197](https://github.com/KYVENetwork/chain/pull/197) Bump to CosmosSDK v0.50.7
- ! (`x/pool`) [#190](https://github.com/KYVENetwork/chain/pull/190) Make inflation-share-weight a decimal.
- [#182](https://github.com/KYVENetwork/chain/pull/182) Make release builds reproducible.
- ! [#183](https://github.com/KYVENetwork/chain/pull/183) Only charge coins which are whitelisted.
- ! (deps) [#174](https://github.com/KYVENetwork/chain/pull/174) Add mainnet KYVE image to interchain tests.
- ! (deps) [#171](https://github.com/KYVENetwork/chain/pull/171) Bump ibc version to v8.2.0 and remove unnecessary dev-dependency.
- ! (deps) [#169](https://github.com/KYVENetwork/chain/pull/169) Migrate to Cosmos-SDK v0.50
- (deps) [#161](https://github.com/KYVENetwork/chain/pull/161) Bump ibc packet forward middleware.

### Bug Fixes

- ! (deps) [#203](https://github.com/KYVENetwork/chain/pull/203) Use newest IAVL version to prevent pruning error.
- [#194](https://github.com/KYVENetwork/chain/pull/194) Return correct result for legacy rpc block_results request.
- [#192](https://github.com/KYVENetwork/chain/pull/192) Cosmos-SDK & IBC swagger-code generation.
- (`x/query`) [#159](https://github.com/KYVENetwork/chain/pull/159) Add pool params back to kyve params query.
- (`x/query`) [#160](https://github.com/KYVENetwork/chain/pull/160) Add funders back to pool query.
- [#163](https://github.com/KYVENetwork/chain/pull/163) Add back swagger UI for KYVE and Cosmos REST queries.

## [v1.4.0](https://github.com/KYVENetwork/chain/releases/tag/v1.4.0) - 2023-11-17

### Features

- ! (`x/funders`) [#141](https://github.com/KYVENetwork/chain/pull/141) Implementation of the new [funders concept](https://commonwealth.im/kyve/discussion/13420-enhancing-kyves-funders-concept).
- [#128](https://github.com/KYVENetwork/chain/pull/128) Add docker build.

### Improvements

- ! (`x/bundles`) [#142](https://github.com/KYVENetwork/chain/pull/142) Halt the pool if a single validator has more than 50% voting power.
- ! (deps) [#33](https://github.com/KYVENetwork/chain/pull/33), [#152](https://github.com/KYVENetwork/chain/pull/152) Upgrade Cosmos SDK to [v0.47.6](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.47.6) ([`v0.47.6-kyve`](https://github.com/KYVENetwork/cosmos-sdk/releases/tag/v0.47.6-kyve-rc0)).

### Bug Fixes

- [#149](https://github.com/KYVENetwork/chain/pull/63) Add bundles query version map to genesis import/export.
- (`x/delegation`) [#135](https://github.com/KYVENetwork/chain/pull/135) CLI: Use correct args amount for tx redelegate.

## [v1.3.2](https://github.com/KYVENetwork/chain/releases/tag/v1.3.2) - 2023-11-16

### Improvements

- (deps) [#138](https://github.com/KYVENetwork/chain/pull/138) Bump ledger dependency to fix ledger build on latest macOS devices.

### Bug Fixes

- (`x/delegation`) [#136](https://github.com/KYVENetwork/chain/pull/136) CLI: Use correct args amount for tx redelegate.
- (`x/bundles`) [#149](https://github.com/KYVENetwork/chain/pull/149) Add bundles query version map to genesis import/export.

## [v1.3.1](https://github.com/KYVENetwork/chain/releases/tag/v1.3.1) - 2023-08-02

### Bug Fixes

- [#122](https://github.com/KYVENetwork/chain/pull/122) Fix makefile go version parse cmd.

## [v1.3.0](https://github.com/KYVENetwork/chain/releases/tag/v1.3.0) - 2023-07-15

### Features

- ! (ibc) [#30](https://github.com/KYVENetwork/chain/pull/30) Integrate [Packet Forward Middleware](https://github.com/strangelove-ventures/packet-forward-middleware).
- ! (`x/bundles`) [#98](https://github.com/KYVENetwork/chain/pull/98) Split inflation rewards between chain and protocol layer.
- ! (`x/bundles`) [#99](https://github.com/KYVENetwork/chain/pull/99) Use weighted round-robin approach for uploader selection.
- ! (`x/bundles`) [#108](https://github.com/KYVENetwork/chain/pull/108) Store stake security for finalized bundles.

### Improvements

- ! (`x/bundles`) [#62](https://github.com/KYVENetwork/chain/pull/62) Payout storage cost directly to the bundle uploader.
- ! (`x/pool`) [#74](https://github.com/KYVENetwork/chain/pull/74) Improve parameter validation in pool proposals.
- ! (`x/stakers`) [#46](https://github.com/KYVENetwork/chain/pull/46) Allow protocol validator commission rewards to be claimed.

### Bug Fixes

- [#96](https://github.com/KYVENetwork/chain/pull/96) Track investor delegation inside auth module.

### Client Breaking

- ! (`x/stakers`) [#46](https://github.com/KYVENetwork/chain/pull/46) Include `MsgClaimCommissionRewards` for claiming commission rewards.

### API Breaking

- (`x/query`) [#87](https://github.com/KYVENetwork/chain/pull/87) Correctly return pools that an account has funded.
- (`x/stakers`) [#46](https://github.com/KYVENetwork/chain/pull/46) Emit an [event](https://github.com/KYVENetwork/chain/blob/release/v1.4.x/x/stakers/spec/05_events.md#eventclaimcommissionrewards) when claiming protocol validator commission rewards.
- (`x/bundles`) [#104](https://github.com/KYVENetwork/chain/pull/104) Improve schema for finalized bundles query.

## [v1.2.3](https://github.com/KYVENetwork/chain/releases/tag/v1.2.3) - 2023-07-15

### API Breaking

- (`x/query`) [#87](https://github.com/KYVENetwork/chain/pull/87) Correctly return pools that an account has funded.
- (`x/bundles`) [#104](https://github.com/KYVENetwork/chain/pull/104) Improve schema for finalized bundles query.

## [v1.2.2](https://github.com/KYVENetwork/chain/releases/tag/v1.2.2) - 2023-06-08

### Bug Fixes

- (deps) [#82](https://github.com/KYVENetwork/chain/pull/82) Bump Cosmos SDK to [v0.46.13](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.46.13) ([`v0.46.13-kyve`](https://github.com/KYVENetwork/cosmos-sdk/releases/tag/v0.46.13-kyve)) to adhere to the [Cosmos SDK Barberry Security Advisory](https://forum.cosmos.network/t/cosmos-sdk-security-advisory-barberry).

## [v1.2.1](https://github.com/KYVENetwork/chain/releases/tag/v1.2.1) - 2023-05-25

### Bug Fixes

- (deps) [#63](https://github.com/KYVENetwork/chain/pull/63) Bump IBC to [v6.1.1](https://github.com/cosmos/ibc-go/releases/tag/v6.1.1) to adhere to the [IBC Huckleberry Security Advisory](https://forum.cosmos.network/t/ibc-security-advisory-huckleberry).

## [v1.2.0](https://github.com/KYVENetwork/chain/releases/tag/v1.2.0) - 2023-05-16

### Bug Fixes

- [#48](https://github.com/KYVENetwork/chain/pull/48) Register Amino types for full Ledger support.
- (`x/team`) [#45](https://github.com/KYVENetwork/chain/pull/45) Adjust vesting schedules of multiple KYVE Core Team members.

## [v1.1.3](https://github.com/KYVENetwork/chain/releases/tag/v1.1.3) - 2023-05-25

### Bug Fixes

- (deps) [#63](https://github.com/KYVENetwork/chain/pull/63) Bump IBC to [v6.1.1](https://github.com/cosmos/ibc-go/releases/tag/v6.1.1) to adhere to the [IBC Huckleberry Security Advisory](https://forum.cosmos.network/t/ibc-security-advisory-huckleberry).

## [v1.1.2](https://github.com/KYVENetwork/chain/releases/tag/v1.1.2) - 2023-05-12

### API Breaking

- (`x/bundles`) [#42](https://github.com/KYVENetwork/chain/pull/42) Emit `VoteEvent` after `BundleProposedEvent` when submitting a bundle.

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
