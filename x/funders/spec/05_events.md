<!--
order: 5
-->

# Events

The funders module contains the following events:

## EventCreateFunder

EventCreateFunder indicates that a new funder has been created.

```protobuf
syntax = "proto3";

message EventCreateFunder {
  // address is the account address of the funder.
  string address = 1;
  // moniker ...
  string moniker = 2;
  // identity is the 64 bit keybase.io identity string
  string identity = 3;
  // website ...
  string website = 4;
  // contact ...
  string contact = 5;
  // description are some additional notes the funder finds important
  string description = 6;
}
```

It gets emitted by the following actions:

- `MsgCreateFunder`

## EventUpdateFunder

EventUpdateFunder indicates that a funder has been updated.

```protobuf
syntax = "proto3";

message EventUpdateFunder {
  // address is the account address of the funder.
  string address = 1;
  // moniker ...
  string moniker = 2;
  // identity is the 64 bit keybase.io identity string
  string identity = 3;
  // website ...
  string website = 4;
  // contact ...
  string contact = 5;
  // description are some additional notes the funder finds important
  string description = 6;
}
```

It gets emitted by the following actions:

- `MsgUpdateFunder`

## EventFundPool

EventFundPool indicates that a funder has provided funds to a pool.

```protobuf
syntax = "proto3";

message EventFundPool {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
  // address is the account address of the pool funder.
  string address = 2;
  // amounts is a list of coins the funder has funded
  repeated cosmos.base.v1beta1.Coin amounts = 3 [
    (gogoproto.nullable)     = false,
    (amino.dont_omitempty)   = true,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
  // amounts_per_bundle is a list of coins the funder wants to distribute per finalized bundle
  repeated cosmos.base.v1beta1.Coin amounts_per_bundle = 4 [
    (gogoproto.nullable)     = false,
    (amino.dont_omitempty)   = true,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
}
```

It gets emitted by the following actions:

- `MsgFundPool`

## EventDefundPool

EventDefundPool indicates that a funder has withdrawn funds from a pool.

```protobuf
syntax = "proto3";

message EventDefundPool {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
  // address is the account address of the pool funder.
  string address = 2;
  // amounts is a list of coins that the funder wants to defund
  repeated cosmos.base.v1beta1.Coin amounts = 3 [
    (gogoproto.nullable)     = false,
    (amino.dont_omitempty)   = true,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
}
```

It gets emitted by the following actions:

- `MsgDefundPool`

## EventPoolOutOfFunds

EventPoolOutOfFunds get emitted when a pool runs out of funds.

```protobuf
syntax = "proto3";

message EventPoolOutOfFunds {
  // pool_id is the unique ID of the pool.
  uint64 pool_id = 1;
}
```
