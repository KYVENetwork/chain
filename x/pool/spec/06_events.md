<!--
order: 6
-->

# Events

The pool module contains the following events:

## EventFundPool

EventFundPool indicates that someone has funded a storage pool with a certain amount.

```protobuf
syntax = "proto3";

message EventFundPool {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
  // address is the account address of the pool funder.
  string address = 2;
  // amount is the amount in ukyve the funder has funded
  uint64 amount = 3;
}
```

It gets emitted by the following actions:

- MsgFundPool

## EventDefundPool

EventDefundPool indicates that someone has defunded a storage pool with a certain amount.

```protobuf
syntax = "proto3";

message EventDefundPool {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
  // address is the account address of the pool funder.
  string address = 2;
  // amount is the amount in ukyve the funder has defunded
  uint64 amount = 3;
}
```

It gets emitted by the following actions:

- MsgDefundPool

## EventPoolFundsSlashed

EventPoolFundsSlashed indicates a funder had not enough KYVE in his funder account anymore to pay for the
validator rewards. In this case the remaining funds get transferred to the treasury and the funder gets
removed.

```protobuf
syntax = "proto3";

message EventPoolFundsSlashed {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
  // address is the account address of the pool funder.
  string address = 2;
  // amount is the amount in ukyve that were transferred to the treasury.
  uint64 amount = 3;
}
```

It gets emitted by the following actions:

- MsgSubmitBundleProposal

## EventPoolOutOfFunds

EventPoolOutOfFunds indicates that a pool has run out of funds and therefore pauses. If that happens someone
has to fund the pool, then the pool will automatically continue.

```protobuf
syntax = "proto3";

message EventPoolOutOfFunds {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
}
```

It gets emitted by the following actions:

- MsgSubmitBundleProposal

## EventCreatePool

EventCreatePool indicates that a new storage pool has been created and is ready to validate and archive.

```protobuf
syntax = "proto3";

message EventCreatePool {
  // id is the unique ID of the pool.
  uint64 id = 1;
  // name is the human readable name of the pool
  string name = 2;
  // runtime is the runtime name of the pool
  string runtime = 3;
  // logo is the logo url of the pool
  string logo = 4;
  // config is either a json stringified config or an
  // external link pointing to the config
  string config = 5;
  // start_key is the first key the pool should start
  // indexing
  string start_key = 6;
  // upload_interval is the interval the pool should validate
  // bundles with
  uint64 upload_interval = 7;
  // inflation_share_weight is the fixed cost which gets paid out
  // to every successful uploader
  uint64 inflation_share_weight = 8;
  // min_delegation is the minimum amount of $KYVE the pool has
  // to have in order to produce bundles
  uint64 min_delegation = 9;
  // max_bundle_size is the max size a data bundle can have
  // (amount of data items)
  uint64 max_bundle_size = 10;
  // version is the current version of the protocol nodes
  string version = 11;
  // binaries points to the current binaries of the protocol node
  string binaries = 12;
  // storage_provider_id is the unique id of the storage provider
  // the pool is archiving the data on
  uint32 storage_provider_id = 13;
  // compression_id is the unique id of the compression type the bundles
  // get compressed with
  uint32 compression_id = 14;
}
```

It gets emitted by the following actions:

- MsgCreatePool


## EventPoolEnabled

EventPoolEnabled indicates that a pool with the given poolId has been enabled.

```protobuf
syntax = "proto3";

message EventPoolEnabled {
  // id is the unique ID of the affected pool.
  uint64 id = 1;
}
```

It gets emitted by the following actions:

- `MsgEnablePool`


## EventPoolDisabled

EventPoolDisabled indicates that a pool with the given poolId has been disabled.

```protobuf
syntax = "proto3";

message EventPoolDisabled {
  // id is the unique ID of the affected pool.
  uint64 id = 1;
}
```

It gets emitted by the following actions:

- `MsgDisablePool`


## EventRuntimeUpgradeScheduled

EventRuntimeUpgradeScheduled indicates that a runtime upgrade has been scheduled.
Pools continue to operate normally until the specified time is reached. 
Then the upgrade is performed.

```protobuf
syntax = "proto3";

message EventRuntimeUpgradeScheduled {
  // runtime is the name of the runtime that will be upgraded.
  string runtime = 1;
  // version is the new version that the runtime will be upgraded to.
  string version = 2;
  // scheduled_at is the time in UNIX seconds when the upgrade will occur.
  uint64 scheduled_at = 3;
  // duration is the amount of seconds the pool will be paused after the
  // scheduled time is reached. This will give node operators time to upgrade
  // their node.
  uint64 duration = 4;
  // binaries contain download links for prebuilt binaries (in JSON format).
  string binaries = 5;
  // affected_pools contains all IDs of pools that will be affected by this runtime upgrade.
  repeated uint64 affected_pools = 6;
}
```

It gets emitted by the following actions:

- `MsgScheduleRuntimeUpgrade`


## EventRuntimeUpgradeCancelled

EventRuntimeUpgradeScheduled indicates that previous scheduled runtime upgrade
was cancelled.

```protobuf
syntax = "proto3";

message EventRuntimeUpgradeCancelled {
  // runtime is the name of the runtime that will be upgraded.
  string runtime = 1;
  // affected_pools contains all IDs of pools that are affected by the
  // cancellation of this runtime upgrade.
  repeated uint64 affected_pools = 2;
}
```

It gets emitted by the following actions:

- `MsgCancelRuntimeUpgrade`


## EventPoolUpdated

EventPoolUpdated indicates that the config of a given pool has been updated. It
emits the raw update string as well as the entire current pool config.

```protobuf
syntax = "proto3";

message EventPoolUpdated {
  // id is the unique ID of the pool.
  uint64 id = 1;
  // raw update string
  string raw_update_string = 2;
  // name is the human readable name of the pool
  string name = 3;
  // runtime is the runtime name of the pool
  string runtime = 4;
  // logo is the logo url of the pool
  string logo = 5;
  // config is either a json stringified config or an
  // external link pointing to the config
  string config = 6;
  // upload_interval is the interval the pool should validate
  // bundles with
  uint64 upload_interval = 7;
  // inflation_share_weight is the fixed cost which gets paid out
  // to every successful uploader
  uint64 inflation_share_weight = 8;
  // min_delegation is the minimum amount of $KYVE the pool has
  // to have in order to produce bundles
  uint64 min_delegation = 9;
  // max_bundle_size is the max size a data bundle can have
  // (amount of data items)
  uint64 max_bundle_size = 10;
  // storage_provider_id is the unique id of the storage provider
  // the pool is archiving the data on
  uint32 storage_provider_id = 11;
  // compression_id is the unique id of the compression type the bundles
  // get compressed with
  uint32 compression_id = 12;
}
```

It gets emitted by the following actions:

- `MsgUpdatePool`
