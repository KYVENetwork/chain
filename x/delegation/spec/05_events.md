<!--
order: 5
-->

# Events

The `x/delegation` module emits the following events:

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

## EventStartUndelegation

EventStartUndelegation is emitted when a user initiates a protocol unbonding.

```protobuf
message EventStartUndelegation {
  // address is the address of the delegator.
  string address = 1;
  // staker is the address of the protocol node.
  string staker = 2;
  // amount is the amount to be undelegated from the protocol node.
  uint64 amount = 3;
  // estimated_undelegation_date is the date in UNIX seconds on when the undelegation will be performed.
  // Note, this number will be incorrect if a governance proposal changes `UnbondingDelegationTime` while unbonding.
  uint64 estimated_undelegation_date = 4;
}
```

## EndBlocker

| Type              | Attribute Key | Attribute Value    |
|-------------------|---------------|--------------------|
| `EventUndelegate` | address       | {delegatorAddress} |
| `EventUndelegate` | staker        | {stakerAddress}    |
| `EventUndelegate` | amount        | {amount}           |

## Messages

### `MsgDelegate`

| Type            | Attribute Key | Attribute Value    |
|-----------------|---------------|--------------------|
| `EventDelegate` | address       | {delegatorAddress} |
| `EventDelegate` | staker        | {stakerAddress}    |
| `EventDelegate` | amount        | {amount}           |

### `MsgRedelegate`

| Type              | Attribute Key | Attribute Value     |
|-------------------|---------------|---------------------|
| `EventRedelegate` | address       | {delegatorAddress}  |
| `EventRedelegate` | from_staker   | {fromStakerAddress} |
| `EventRedelegate` | to_staker     | {toStakerAddress}   |
| `EventRedelegate` | amount        | {amount}            |

### `MsgWithdrawRewards`

| Type                   | Attribute Key | Attribute Value    |
|------------------------|---------------|--------------------|
| `EventWithdrawRewards` | address       | {delegatorAddress} |
| `EventWithdrawRewards` | staker        | {stakerAddress}    |
| `EventWithdrawRewards` | amounts       | {amounts}          |
