<!--
order: 5
-->

# Events

The `x/stakers` module contains the following events:

## EventUpdateParams

EventUpdateParams is emitted when the parameters were changed by the governance.

```protobuf
message EventUpdateParams {
  // old_params is the module's old parameters.
  kyve.bundles.v1beta1.Params old_params = 1 [(gogoproto.nullable) = false];
  // new_params is the module's new parameters.
  kyve.bundles.v1beta1.Params new_params = 2 [(gogoproto.nullable) = false];
  // payload is the parameter updates that were performed.
  string payload = 3;
}
```

## EventCreateStaker

EventBundleProposed indicates that a new staker was created.

```protobuf
message EventCreateStaker {
  // staker is the account address of the protocol node.
  string staker = 1;
  // amount for inital self-delegation
  uint64 amount = 2;
  // commission
  string commission = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```

It gets thrown from the following actions:

- MsgCreateStaker

## EventUpdateMetadata

EventUpdateMetadata is an event emitted when a protocol node updates their
metadata.

```protobuf
message EventUpdateMetadata {
  // staker is the account address of the protocol node.
  string staker = 1;
  // moniker ...
  string moniker = 2;
  // website ...
  string website = 3;
  // identity ...
  string identity = 4;
  // security_contact ...
  string security_contact = 5;
  // details ...
  string details = 6;
}
```

It gets thrown from the following actions:

- MsgUpdateMetadata

## EventUpdateCommission

EventUpdateCommission indicates that a staker has changes its commission.

```protobuf
message EventUpdateCommission {
  // staker is the account address of the protocol node.
  string staker = 1;
  // commission ...
  string commission = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```

It gets thrown from the following actions:

- EndBlock

## EventJoinPool

EventClaimUploaderRole indicates that a staker has joined a pool.

```protobuf
message EventJoinPool {
  // pool_id is the pool the staker joined
  uint64 pool_id = 1;
  // staker is the address of the staker
  string staker = 2;
  // valaddress is the address of the protocol node which 
  // votes in favor of the staker
  string valaddress = 3;
  // amount is the amount of funds transferred to the valaddress
  uint64 amount = 4;
}
```

It gets thrown from the following actions:

- MsgJoinPool

## EventLeavePool

EventLeavePool indicates that a staker has left a pool.
Either by leaving or by getting kicked out for the following reasons:

- misbehaviour (usually together with a slash)
- all pool slots are taken and a node with more stake joined.

```protobuf
message EventLeavePool {
  // pool_id ...
  uint64 pool_id = 1;
  // staker ...
  string staker = 2;
}
```

It gets thrown from the following actions:

- EndBlock
- bundles/MsgSubmitBundleProposal
- MsgJoinPool