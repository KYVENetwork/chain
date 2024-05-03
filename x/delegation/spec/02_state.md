<!--
order: 2
-->

# State

The module is mainly responsible for handling the f1-distribution state.
Furthermore, it is also responsible for unbonding and redelegation. 

## F1-Distribution
The state is split across four proto-files which all have their own
prefix in the KV-Store.

### DelegationData
DelegationData exist for every staker and stores primarily 
total delegation, the rewards for the current period and keeps track
of the f1-index. It exists as long as the staker has at least `1ukyve` delegation.

- DelegationData: `0x03 | StakerAddr -> ProtocolBuffer(stakerDelegationData)`

```go
type DelegationData struct {
    // Every staker has one DelegationData
    Staker string
    CurrentRewards github_com_cosmos_cosmos_sdk_types.Coins
    TotalDelegation uint64
    LatestIndexK uint64
    // delegator_count the amount of different addresses delegating to the staker
    DelegatorCount uint64
    // latest_index_was_undelegation helps indicates when an entry can be deleted
    LatestIndexWasUndelegation bool
}
```

### Delegator
Delegator represents a pair of (staker, delegator) and the corresponding f1-index.

- Delegator: `0x01 | 0x00 | StakerAddr | DelegatorAddr -> ProtocolBuffer(delegator)` 

One additional index is maintained to query for all stakers a delegator has delegated to:

- DelegatorIndex2: `0x01 | 0x01 | DelegatorAddr | StakerAddr -> ProtocolBuffer(delegator)`

```go
type Delegator struct {
    Staker string
    Delegator string
    KIndex uint64
    InitialAmount uint64
}
```

### DelegationEntry
DelegationEntries are used internally by the f1-distribution.
They mark the beginning of every period.

- DelegationEntry: `0x02 | StakerAddr | kIndex -> ProtocolBuffer(delegationEntry)`

```go
type DelegationEntry struct {
    Staker string
    KIndex uint64
    Value sdk.Dec
}
```

### DelegationSlash
DelegationSlash represents an internal f1-slash.
It is needed to calculate the actual amount of stake
after a slash occurred.

- DelegationSlash: `0x04 | StakerAddr | kIndex -> ProtocolBuffer(delegationSlash)`

```go
type DelegationSlash struct {
    Staker string
    KIndex uint64
    Fraction sdk.Dec
}
```

## Unbonding Queue

### QueueState
For the unbonding queue the app needs to keep track of the head (HighIndex) and
tail (LowIndex) of the queue. New entries are appended to the
head. The EndBlocker checks the tail if entries are due and processes them.

- QueueState: `0x05 -> ProtocolBuffer(queueState)`

```go
type DelegationSlash struct {
    LowIndex uint64
    HighIndex uint64
}
```

### UndelegationQueueEntry
Every time a user starts an undelegation an entry is created 
and appended to the head of the queue. I.e. the current HighIndex is
incremented and assigned to the entry.
The order of the queue is automatically provided by the KV-Store.

- UndelegationQueueEntry: `0x06 | 0x00 | Index  -> ProtocolBuffer(undelegationQueueEntry)`

A second index is provided so that users can query their own pending entries
without iterating the entire queue.

- UndelegationQueueEntryIndex2: `0x06 | 0x01 | StakerAddr | Index  -> ProtocolBuffer(undelegationQueueEntry)`


```go
type UndelegationQueueEntry struct {
    Index uint64
	Staker string
	Delegator string
    Amount uint64
    CreationTime uint64
}
```


## Redelegation Spells

Redelegation spells do not require a queue for tracking expired
spells, as they are checked on demand when the users trys to
redelegate. 

### RedelegationCooldown

Every used redelegation spell is stored in the KV-Store with its creation time.
Once the oldest entry is older then `RedelegationCooldown` it can be reused.
To avoid keeping track of a global counter we use the blockHeight to generate
a unique key for the KV-Store. 
Therefore, it is only possible to perform one redelegation per block.

- RedelegationCooldown: `0x07 | DelegatorAddr | blockHeight -> ProtocolBuffer(redelegationCooldown)`

```go
type RedelegationCooldown struct {
    Address String
    CreationDate uint64
}
```

