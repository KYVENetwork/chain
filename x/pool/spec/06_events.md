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
  // amount is the amount in ukyve that were transferred to the treasury
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
  // operating_cost is the fixed cost which gets paid out
  // to every successful uploader
  uint64 operating_cost = 8;
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
