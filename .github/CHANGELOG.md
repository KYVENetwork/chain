<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking Protobuf, gRPC and REST routes used by end-users.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [v0.4.0](https://github.com/KYVENetwork/chain/releases/tag/v0.4.0) - 2022-06-7

### Features

- Implemented scheduled upgrades for pool versions
- Implemented `abstain` vote besides `valid` and `invalid`. Validators who don't vote 5 times in a row at all get removed with a timeout slash

### Client Breaking Changes

- The arg `vote` on `MsgVoteProposal` changed from `bool` to `uint64`. 0 = valid, 1 = invalid, 2 = abstain
- The arg `versions` on `MsgCreatePoolProposal` changed to `version`
- The arg `binaries` got added to `MsgCreatePoolProposal`

### Improvements

- Check the quorum of the bundle proposal on chain to prevent unjustified slashes
- Don't drop bundle proposals if one funder can't afford the funding cost, instead remove all of them and proceed
- If a validator submits a `NO_DATA_BUNDLE` the will just skip the upload instead of proposing an empty bundle
- Added query `QueryFunder`
- Added query `QueryStaker`
- Added query `QueryDelegator`

### Bug Fixes

### Deprecated

- Deprecated `versions` on `kyve.registry.v1beta1.Pool`
