# CHANGELOG

## v0.7.0_beta7

- refactor: refactored custom keys by renaming `height` to `index` and adding both properties `from_key` and `to_key` to bundle proposal
- refactor: renamed `current_height` to `current_index` and `current_value` to `current_summary` on pool
- refactor: removed `from_height` and `to_height` from bundle proposal and instead added `from_index` and `bundle_size` to indicate more clearly the data range of the bundle
- refactor: renamed `bundle_hash` to `data_hash` on bundle proposal to make it clear the raw compressed data as it lies on the storage provider is hashed
- refactor: renamed `byte_size` to `data_size` on bundle proposal
- refactor: refactored bundle value by renaming `to_value` to `bundle_summary` and allowing protocol nodes to submit an entire bundle summary on-chain instead of just a single value
- feat: added and implemented event `EventPointIncreased`
- feat: added and implemented event `EventPointsReset`
- fix: implemented unused event `EventSlash`
- fix: throw error now if staker joins with a valaddress that is already used by another staker in the same pool

## v0.7.0_beta8
- refactor: added `ar://` to every arweave tx for pool logos
- feat: pool config is now stored externally on arweave of ipfs
- feat: `storageProviderId` and `compressionId` were introduced to pools to enable dynamic storage provider and compression switching
- Refactor Events:
  - Emit ClaimedUploaderRole-event
  - EventDelegate: `node` -> `staker`
  - EventUndelegate: `node` -> `staker`
  - EventRedelegate: `from_node` -> `from_staker`, `to_node` -> `to_staker`
  - EventWithdrawRewards: `from_node` -> `staker`
  - EventCreateStaker: `address` -> `staker`
  - EventUpdateMetadata: `address` -> `staker`
  - EventSlash: `address` -> `staker`
  - EventUpdateCommission: `address` -> `staker`
- Emit `LeavePoolEvent` if staker gets kicked out of pool
