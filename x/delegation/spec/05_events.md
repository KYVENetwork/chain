<!--
order: 5
-->

# Events

The `x/delegation` module emits the following events:

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
| `EventWithdrawRewards` | amount        | {amount}           |
