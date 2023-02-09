<!--
order: 6
-->

# Events

The bundles module contains the following events:

## EventBundleProposed

EventBundleProposed indicates that a new bundle proposal was submitted
to a storage pool. This event contains all information about the
proposal.

```protobuf
syntax = "proto3";

message EventBundleProposed {
    // pool_id is the unique ID of the pool.
    uint64 pool_id = 1;
    // internal id for the KYVE-bundle
    uint64 id = 2;
    // storage_id is the ID to retrieve to data item from the configured storage provider
    // e.g. the ARWEAVE-id
    string storage_id = 3;
    // Address of the uploader/proposer of the bundle
    string uploader = 4;
    // data_size size in bytes of the data
    uint64 data_size = 5;
    // from_index starting index of the bundle (inclusive)
    uint64 from_index = 6;
    // bundle_size amount of data items in the bundle
    uint64 bundle_size = 7;
    // from_key the key of the first data item in the bundle
    string from_key = 8;
    // to_key the key of the last data item in the bundle
    string to_key = 9;
    // bundle_summary is a short string holding some useful information of
    // the bundle which will get stored on-chain
    string bundle_summary = 10;
    // data_hash is a sha256 hash of the raw compressed data
    string data_hash = 11;
    // proposed_at the unix time when the bundle was proposed
    uint64 proposed_at = 12;
    // storage_provider_id the unique id of the storage provider where
    // the data of the bundle is tored
    uint32 storage_provider_id = 13;
    // compression_id  the unique id of the compression type the data
    // of the bundle was compressed with
    uint32 compression_id = 14;
}
```

It gets thrown from the following actions:

- MsgSubmitBundleProposal

## EventBundleFinalized

EventBundleFinalized indicates that a bundle has been finalized with
a certain vestingStatus. This vestingStatus includes dropped, invalid and valid
bundles.

```protobuf
syntax = "proto3";

message EventBundleFinalized {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
  // internal id for the KYVE-bundle
  uint64 id = 2;
  // total voting power which voted for valid
  uint64 valid = 3;
  // total voting power which voted for invalid
  uint64 invalid = 4;
  // total voting power which voted for abstain
  uint64 abstain = 5;
  // total voting power of the pool
  uint64 total = 6;
  // vestingStatus of the finalized bundle
  BundleStatus vestingStatus = 7;
  // rewards transferred to treasury (in ukyve)
  uint64 reward_treasury = 8;
  // rewardUploader rewards directly transferred to uploader (in ukyve)
  uint64 reward_uploader = 9;
  // rewardDelegation rewards distributed among all delegators (in ukyve)
  uint64 reward_delegation = 10;
  // rewardTotal the total bundle reward
  uint64 reward_total = 11;
  // finalized_at the block height where the bundle got finalized
  uint64 finalized_at = 12;
  // uploader the address of the uploader of this bundle
  string uploader = 13;
  // next_uploader the address of the next uploader after this bundle
  string next_uploader = 14;
}
```

It gets thrown from the following actions:

- MsgSubmitBundleProposal
- EndBlock

## EventBundleVote

EventBundleVote indicates that a participant has voted on a bundle.

```protobuf
syntax = "proto3";

message EventBundleVote {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
  // staker is the account staker of the protocol node.
  string staker = 2;
  // storage_id is the unique ID of the bundle.
  string storage_id = 3;
  // vote is for what the validator voted with
  VoteType vote = 4;
}
```

It gets thrown from the following actions:

- MsgVoteBundleProposal

## EventClaimUploaderRole

EventClaimUploaderRole indicates that a participant has claimed
a free uploader role spot in a storage pool.

```protobuf
syntax = "proto3";

message EventClaimedUploaderRole {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
  // id internal id for the KYVE-bundle
  uint64 id = 2;
  // new_uploader the address of the participant who claimed
  // the free uploader role
  string new_uploader = 3;
}
```

It gets thrown from the following actions:

- MsgClaimUploaderRole

## EventSkippedUploaderRole

EventSkippedUploaderRole indicates that the current uploader
of a storage pool has skipped his uploader role.

```protobuf
syntax = "proto3";

message EventClaimedUploaderRole {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
  // id internal id for the KYVE-bundle
  uint64 id = 2;
  // new_uploader the address of the participant who claimed
  // the free uploader role
  string new_uploader = 3;
}
```

It gets thrown from the following actions:

- MsgSkipUploaderRole

## EventPointIncreased

EventPointIncreased indicates that a staker received a point
for being offline while voting or submitting.

```protobuf
syntax = "proto3";

message EventPointIncreased {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
  // staker is the address of the staker who received the point
  string staker = 2;
  // current_points is the amount of points the staker has now
  uint64 current_points = 3;
}
```

It gets thrown from the following actions:

- MsgSubmitBundleProposal
- EndBlock

## EventPointsReset

EventPointsReset indicates that a staker who previously had
some points got his points reset due to being active again.

```protobuf
syntax = "proto3";

message EventPointsReset {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
  // staker is the address of the staker who has zero points now
  string staker = 2;
}
```

It gets thrown from the following actions:

- MsgSubmitBundleProposal
- MsgVoteBundleProposal
- MsgSkipUploaderRole

